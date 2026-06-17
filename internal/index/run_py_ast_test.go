package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestRunPyAst_Success runs the real ast script against a valid file and parses
// the JSON it dumps. Skipped without a Python interpreter.
func TestRunPyAst_Success(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "m.py")
	if err := os.WriteFile(f, []byte("def foo():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runPyAst(python, pyAstDefScript, f)
	if err != nil {
		t.Fatalf("runPyAst: %v", err)
	}
	var raw []pyAstFunc
	if err := json.Unmarshal(out, &raw); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, out)
	}
	if len(raw) != 1 || raw[0].Name != "foo" {
		t.Fatalf("ast dump = %+v, want one func foo", raw)
	}
}

// TestRunPyAst_SyntaxError proves a SyntaxError (script exit 2) surfaces as an
// error so the caller per-file-falls-back to the line indexer.
func TestRunPyAst_SyntaxError(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "broken.py")
	if err := os.WriteFile(f, []byte("def (:\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runPyAst(python, pyAstDefScript, f); err == nil {
		t.Fatal("expected error for a file with a SyntaxError")
	}
}

// TestRunPyAst_InterpreterMissing proves an unrunnable interpreter is reported as
// an error (the Run() failure branch).
func TestRunPyAst_InterpreterMissing(t *testing.T) {
	if _, err := runPyAst("definitely-not-a-real-interpreter-xyz", pyAstDefScript, "x.py"); err == nil {
		t.Fatal("expected error for a missing interpreter")
	}
}
