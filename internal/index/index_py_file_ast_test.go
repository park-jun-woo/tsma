package index

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIndexPyFileast_Success exercises the per-file precise step end to end (ast
// subprocess + parse). Skipped without a Python interpreter.
func TestIndexPyFileast_Success(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}
	dir := t.TempDir()
	abs := filepath.Join(dir, "calc.py")
	if err := os.WriteFile(abs, []byte("def classify(n):\n    return n\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := indexPyFileast("src/calc.py", abs, python)
	if err != nil {
		t.Fatalf("indexPyFileast: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "classify" {
		t.Fatalf("funcs = %+v, want one classify", funcs)
	}
	if funcs[0].QualifiedName != "src.classify" {
		t.Errorf("qualified = %q, want src.classify", funcs[0].QualifiedName)
	}
}

// TestIndexPyFileast_SubprocessError proves a subprocess failure is returned as
// an error (the err != nil branch), so Index falls back to the line indexer.
func TestIndexPyFileast_SubprocessError(t *testing.T) {
	if _, err := indexPyFileast("x.py", "x.py", "definitely-not-a-real-interpreter-xyz"); err == nil {
		t.Fatal("expected error when the interpreter is unrunnable")
	}
}
