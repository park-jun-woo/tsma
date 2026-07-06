package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeTSXML is a canned TS parse tree containing one bare call (targetFn())
// and one construction (new TargetClass()) — the two node kinds
// collectTSCalledNames harvests.
const fakeTSXML = `<program srow="0" scol="0" erow="4" ecol="0">
 <call_expression srow="1" scol="0" erow="1" ecol="10">
  <identifier field="function" srow="1" scol="0" erow="1" ecol="8">targetFn</identifier>
 </call_expression>
 <new_expression srow="2" scol="0" erow="2" ecol="20">
  <identifier field="constructor" srow="2" scol="4" erow="2" ecol="15">TargetClass</identifier>
 </new_expression>
</program>`

// writeTSFakePackage lays out a minimal TS package (pkg/math.ts with a sibling
// math.test.ts and a __tests__/extra.spec.ts) under a fresh root and returns it.
func writeTSFakePackage(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range map[string]string{
		"pkg/math.ts":                 "export const x = 1;\n",
		"pkg/math.test.ts":            "test\n",
		"pkg/__tests__/extra.spec.ts": "spec\n",
	} {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// TestBuildTSPkgTestIndexFakeCLI covers the CLI-present branches without a real
// tree-sitter install: test files in both the package dir and its __tests__/
// are ingested, and a package with no test files yields nil (len==0).
func TestBuildTSPkgTestIndexFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_TS_GRAMMAR", fakeTSXML)
	root := writeTSFakePackage(t)

	idx := BuildTSPkgTestIndex(root, "pkg")
	if idx == nil {
		t.Fatal("expected a populated index for pkg")
	}
	got := idx.refs["targetFn"]
	want := []string{
		filepath.Join("pkg", "math.test.ts"),
		filepath.Join("pkg", "__tests__", "extra.spec.ts"),
	}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("refs[targetFn] = %v, want %v", got, want)
	}
	if got := idx.refs["TargetClass"]; len(got) != 2 {
		t.Errorf("refs[TargetClass] = %v, want both test files", got)
	}

	// A package dir with no test files (and no __tests__/) -> nil.
	if err := os.MkdirAll(filepath.Join(root, "lonely"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := BuildTSPkgTestIndex(root, "lonely"); got != nil {
		t.Errorf("expected nil index for a dir without test files, got %+v", got)
	}
}

// TestIngestTSDirFakeCLI covers the per-entry branches deterministically: a
// subdirectory (IsDir skip), a non-test .ts file (isTSTestFile skip), and a
// *.test.ts file (ingested).
func TestIngestTSDirFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeTSXML)
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"plain.ts", "foo.test.ts"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	idx := &TSPkgTestIndex{refs: map[string][]string{}}
	ingestTSDir(idx, dir, ".", script, "")
	if got := idx.refs["targetFn"]; len(got) != 1 || got[0] != "foo.test.ts" {
		t.Errorf("refs[targetFn] = %v, want [foo.test.ts] (only the test file ingested)", got)
	}
}

// TestIngestTSTestFileFakeCLI covers the successful parse + ref recording
// branch without a real CLI.
func TestIngestTSTestFileFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeTSXML)
	idx := &TSPkgTestIndex{refs: map[string][]string{}}
	ingestTSTestFile(idx, script, "", "/x/foo.test.ts", "foo.test.ts")
	if got := idx.refs["targetFn"]; len(got) != 1 || got[0] != "foo.test.ts" {
		t.Errorf("refs[targetFn] = %v, want [foo.test.ts]", got)
	}
}

// TestTypeScriptFuncMatcherMatchFuncFakeCLI covers the content-index branches
// of MatchFunc: a referenced name attributes to the calling test file, and an
// unreferenced name falls through to the filename fallback.
func TestTypeScriptFuncMatcherMatchFuncFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_TS_GRAMMAR", fakeTSXML)
	root := writeTSFakePackage(t)
	m := &TypeScriptFuncMatcher{}

	// Content hit: targetFn is called by both test files.
	fn := &model.Function{Name: "targetFn", File: "pkg/math.ts"}
	tm, ok := m.MatchFunc(root, fn)
	if !ok || len(tm.Files) != 2 || tm.Files[0] != filepath.Join("pkg", "math.test.ts") {
		t.Errorf("content hit: ok=%v tm=%+v", ok, tm)
	}

	// Content miss with a populated index -> filename fallback finds the
	// conventional sibling math.test.ts.
	miss := &model.Function{Name: "absent", File: "pkg/math.ts"}
	tm, ok = m.MatchFunc(root, miss)
	if !ok || len(tm.Files) != 1 || tm.Files[0] != filepath.Join("pkg", "math.test.ts") {
		t.Errorf("content miss fallback: ok=%v tm=%+v", ok, tm)
	}
}

// TestTypeScriptFuncMatcherMatchFuncNoCLI covers the idx==nil branch: with
// tree-sitter absent MatchFunc degrades to the filename fallback (found and
// not found).
func TestTypeScriptFuncMatcherMatchFuncNoCLI(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	root := writeTSFakePackage(t)
	m := &TypeScriptFuncMatcher{}

	fn := &model.Function{Name: "targetFn", File: "pkg/math.ts"}
	tm, ok := m.MatchFunc(root, fn)
	if !ok || len(tm.Files) != 1 || tm.Files[0] != filepath.Join("pkg", "math.test.ts") {
		t.Errorf("no-CLI fallback: ok=%v tm=%+v", ok, tm)
	}

	if err := os.MkdirAll(filepath.Join(root, "bare"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bare", "lonely.ts"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	miss := &model.Function{Name: "absent", File: "bare/lonely.ts"}
	if _, ok := m.MatchFunc(root, miss); ok {
		t.Error("no-CLI fallback without a conventional test file should report not found")
	}
}
