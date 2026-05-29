package coverage

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFindCargo(t *testing.T) {
	path, err := findCargo()
	if _, lookErr := exec.LookPath("cargo"); lookErr == nil {
		if err != nil {
			t.Errorf("cargo on PATH but findCargo returned error: %v", err)
		}
		if path == "" {
			t.Error("expected non-empty cargo path")
		}
		return
	}
	if err == nil {
		t.Error("expected error when cargo is not installed")
	}
}

func TestFindCargo_found(t *testing.T) {
	// Put a fake executable named "cargo" on PATH so the success branch runs
	// deterministically regardless of whether the real toolchain is installed.
	dir := t.TempDir()
	fake := filepath.Join(dir, "cargo")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	path, err := findCargo()
	if err != nil {
		t.Fatalf("expected cargo to be found, got error: %v", err)
	}
	if path != fake {
		t.Errorf("expected path %q, got %q", fake, path)
	}
}

func TestFindCargo_notFound(t *testing.T) {
	// Empty PATH ensures cargo cannot be located -> error branch.
	t.Setenv("PATH", t.TempDir())
	_, err := findCargo()
	if err == nil {
		t.Fatal("expected error when cargo is not on PATH")
	}
}
