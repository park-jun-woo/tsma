package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// writeFakeExec creates an executable shell stub named bin in dir.
func writeFakeExec(t *testing.T, dir, bin string) {
	t.Helper()
	p := filepath.Join(dir, bin)
	if err := os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

// TestResolvePython covers all three interpreter-resolution branches by
// manipulating PATH: python3 preferred, python as fallback, and "" when neither
// is present (the signal that triggers the line-based fallback).
func TestResolvePython(t *testing.T) {
	t.Run("python3 preferred when present", func(t *testing.T) {
		if _, err := exec.LookPath("python3"); err != nil {
			t.Skip("python3 not on PATH")
		}
		if got := resolvePython(); got != "python3" {
			t.Errorf("resolvePython() = %q, want python3", got)
		}
	})

	t.Run("falls back to python when only python present", func(t *testing.T) {
		dir := t.TempDir()
		writeFakeExec(t, dir, "python")
		t.Setenv("PATH", dir)
		if got := resolvePython(); got != "python" {
			t.Errorf("resolvePython() = %q, want python", got)
		}
	})

	t.Run("empty when no interpreter on PATH", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if got := resolvePython(); got != "" {
			t.Errorf("resolvePython() = %q, want empty", got)
		}
	})
}
