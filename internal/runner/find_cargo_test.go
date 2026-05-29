package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFindCargo(t *testing.T) {
	path, err := findCargo()
	if _, lookErr := exec.LookPath("cargo"); lookErr == nil {
		// cargo present: expect a path and no error.
		if err != nil {
			t.Errorf("cargo on PATH but findCargo returned error: %v", err)
		}
		if path == "" {
			t.Error("expected non-empty cargo path")
		}
		return
	}
	// cargo absent: expect a clear error.
	if err == nil {
		t.Error("expected error when cargo is not installed")
	}
}

func TestFindCargo_found(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "cargo")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	path, err := findCargo()
	if err != nil {
		t.Fatalf("expected cargo found, got error: %v", err)
	}
	if path != fake {
		t.Errorf("expected %q, got %q", fake, path)
	}
}

func TestFindCargo_notFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := findCargo(); err == nil {
		t.Fatal("expected error when cargo not on PATH")
	}
}
