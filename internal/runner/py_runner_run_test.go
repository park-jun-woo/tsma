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
