package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCoveragePython(t *testing.T) {
	result := findCoveragePython()
	// Should return either "python3" or "python"
	if result != "python3" && result != "python" {
		t.Errorf("findCoveragePython() = %q, want 'python3' or 'python'", result)
	}
}

func TestFindCoveragePython_python3Present(t *testing.T) {
	// Provide a fake python3 on PATH so the python3 branch is taken.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "python3"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	if got := findCoveragePython(); got != "python3" {
		t.Errorf("expected python3, got %q", got)
	}
}

func TestFindCoveragePython_fallback(t *testing.T) {
	// Empty PATH: python3 not found -> fallback to "python".
	t.Setenv("PATH", t.TempDir())
	if got := findCoveragePython(); got != "python" {
		t.Errorf("expected fallback 'python', got %q", got)
	}
}
