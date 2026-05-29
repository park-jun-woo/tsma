package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPyRunnerDetectsFramework(t *testing.T) {
	// Verify that detectPytest is used in the Run logic:
	// project with conftest.py should use pytest
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "conftest.py"), []byte("# conftest"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectPytest(dir) {
		t.Error("detectPytest should return true for project with conftest.py")
	}

	// project without pytest indicators should use unittest
	emptyDir := t.TempDir()
	if detectPytest(emptyDir) {
		t.Error("detectPytest should return false for empty project")
	}
}

func TestPyRunnerFindPython(t *testing.T) {
	python := findPython()
	if python != "python3" && python != "python" {
		t.Errorf("findPython() = %q, want python3 or python", python)
	}
}

// fakePython installs a fake python3 on PATH (prepended) with the given body.
func fakePython(t *testing.T, body string) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "python3"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestPyRunnerRun_pytestBranchPass(t *testing.T) {
	dir := t.TempDir()
	// conftest.py -> detectPytest true -> pytest branch.
	if err := os.WriteFile(filepath.Join(dir, "conftest.py"), []byte("# conftest"), 0o644); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "test_x.py"), []byte("def test_x():\n    assert True\n"), 0o644)
	fakePython(t, "exit 0\n")

	r := &PyRunner{}
	res, err := r.Run(dir, "test_x.py")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestPyRunnerRun_unittestBranchFail(t *testing.T) {
	dir := t.TempDir()
	// No pytest indicators -> unittest branch.
	os.WriteFile(filepath.Join(dir, "test_x.py"), []byte("# noop\n"), 0o644)
	fakePython(t, "exit 1\n")

	r := &PyRunner{}
	res, err := r.Run(dir, "test_x.py")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false when python exits non-zero")
	}
}
