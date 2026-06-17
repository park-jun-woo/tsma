package index

import (
	"os"
	"path/filepath"
	"testing"
)

// findFn returns the first indexed function with the given bare name.
func findFn(funcs []fnLite, name string) (fnLite, bool) {
	for _, f := range funcs {
		if f.name == name {
			return f, true
		}
	}
	return fnLite{}, false
}

type fnLite struct {
	name      string
	receiver  string
	startLine int
	endLine   int
}

// TestPyAstIndexer_PrecisePython exercises the real ast subprocess path. It is
// skipped when no Python interpreter is on PATH (clean environments fall back to
// the line indexer, covered separately). It asserts ast-only precision: methods
// carry their class receiver, a function nested inside a method does NOT, async
// defs are found, and line ranges are exact.
func TestPyAstIndexer_PrecisePython(t *testing.T) {
	if resolvePython() == "" {
		t.Skip("no python interpreter on PATH; ast precise path unavailable")
	}
	dir := t.TempDir()
	src := `def classify(n):
    if n > 0:
        return "pos"
    return "nonpos"


async def fetch(value):
    return value


class Calculator:
    def add(self, a, b):
        def _inner(v):
            return v + 1
        return _inner(a) + b

    @staticmethod
    def mul(a, b):
        return a * b
`
	if err := os.WriteFile(filepath.Join(dir, "calc.py"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyAstIndexer{python: resolvePython(), fallback: &PyIndexer{}}
	modelFuncs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyAstIndexer.Index: %v", err)
	}

	var funcs []fnLite
	for _, f := range modelFuncs {
		funcs = append(funcs, fnLite{f.Name, f.Receiver, f.StartLine, f.EndLine})
	}

	for _, want := range []string{"classify", "fetch", "add", "_inner", "mul"} {
		if _, ok := findFn(funcs, want); !ok {
			t.Errorf("ast index missing function %q; got %+v", want, funcs)
		}
	}

	add, _ := findFn(funcs, "add")
	if add.receiver != "Calculator" {
		t.Errorf("method add receiver = %q, want Calculator", add.receiver)
	}
	if add.startLine != 12 || add.endLine != 15 {
		t.Errorf("method add range = %d-%d, want 12-15", add.startLine, add.endLine)
	}

	inner, _ := findFn(funcs, "_inner")
	if inner.receiver != "" {
		t.Errorf("nested function _inner receiver = %q, want \"\" (not a method)", inner.receiver)
	}

	classify, _ := findFn(funcs, "classify")
	if classify.startLine != 1 || classify.endLine != 4 {
		t.Errorf("classify range = %d-%d, want 1-4", classify.startLine, classify.endLine)
	}
}

// TestPyAstIndexer_FallbackWhenNoPython proves graceful fallback: a PyAstIndexer
// with no interpreter delegates to the line-based PyIndexer (zero regression).
func TestPyAstIndexer_FallbackWhenNoPython(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.py"), []byte("def foo():\n    pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := &PyAstIndexer{python: "", fallback: &PyIndexer{}}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	found := false
	for _, f := range funcs {
		if f.Name == "foo" {
			found = true
		}
	}
	if !found {
		t.Fatalf("fallback line indexer did not find foo; got %+v", funcs)
	}
}

// TestPyAstIndexer_IndexBranches drives the Index walk through every filter
// branch in one tree: a .tsmignore-matched directory (SkipDir) and file (skip),
// a skipPyDir directory (.venv), a non-Python file (isPySource false), a
// syntax-error file (ast fails → per-file line fallback), and a clean file (ast
// path). Skipped without a Python interpreter (the precise path needs one).
func TestPyAstIndexer_IndexBranches(t *testing.T) {
	if resolvePython() == "" {
		t.Skip("no python interpreter on PATH; precise path unavailable")
	}
	dir := t.TempDir()

	write := func(rel, content string) {
		abs := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(".tsmignore", "ignored_dir/\nignore_me.py\n")
	write("ignored_dir/buried.py", "def buried():\n    pass\n")
	write("ignore_me.py", "def ignored_top():\n    pass\n")
	write(".venv/lib.py", "def vendored():\n    pass\n")
	write("notes.txt", "not python\n")
	write("broken.py", "def (:\n") // SyntaxError → ast exit 2 → line fallback
	write("good.py", "def classify(n):\n    return n\n")

	idx := &PyAstIndexer{python: resolvePython(), fallback: &PyIndexer{}}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	names := map[string]bool{}
	for _, f := range funcs {
		names[f.Name] = true
	}
	if !names["classify"] {
		t.Errorf("expected ast path to index classify; got %v", names)
	}
	// broken.py triggers the per-file ast failure → line-based fallback branch;
	// the indexer must not crash and must still index the clean file.
	for _, gone := range []string{"buried", "ignored_top", "vendored"} {
		if names[gone] {
			t.Errorf("function %q should have been filtered out; got %v", gone, names)
		}
	}
}

// TestPyAstIndexer_IndexWalkError covers the walk callback's err != nil branch:
// an unreadable directory makes filepath.Walk invoke the callback with a
// non-nil error, which Index swallows (return nil) so one bad directory never
// aborts the whole scan. Skipped without a Python interpreter or when running as
// root (root bypasses the permission that creates the error).
func TestPyAstIndexer_IndexWalkError(t *testing.T) {
	if resolvePython() == "" {
		t.Skip("no python interpreter on PATH; precise path unavailable")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root; directory permissions are bypassed")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "good.py"), []byte("def ok():\n    pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	denied := filepath.Join(dir, "denied")
	if err := os.Mkdir(denied, 0o000); err != nil {
		t.Fatal(err)
	}
	// Restore perms so t.TempDir cleanup can remove it.
	t.Cleanup(func() { _ = os.Chmod(denied, 0o755) })

	idx := &PyAstIndexer{python: resolvePython(), fallback: &PyIndexer{}}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("Index must swallow the walk error, got: %v", err)
	}
	found := false
	for _, f := range funcs {
		if f.Name == "ok" {
			found = true
		}
	}
	if !found {
		t.Errorf("the readable file should still be indexed despite the denied dir; got %+v", funcs)
	}
}
