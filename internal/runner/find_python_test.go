package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindPythonReturnsValidBinary(t *testing.T) {
	result := findPython()
	if result != "python3" && result != "python" {
		t.Errorf("findPython() = %q, want \"python3\" or \"python\"", result)
	}
}

func TestFindPython_python3Present(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "python3"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	if got := findPython(); got != "python3" {
		t.Errorf("expected python3, got %q", got)
	}
}

func TestFindPython_fallback(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if got := findPython(); got != "python" {
		t.Errorf("expected fallback python, got %q", got)
	}
}
