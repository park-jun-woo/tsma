package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRunCargoLLVMCovBadBinary verifies the command wrapper surfaces an error
// when the cargo binary cannot be executed. A successful run requires a real
// cargo + cargo-llvm-cov toolchain (E2E only).
func TestRunCargoLLVMCovBadBinary(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "llvm-cov.json")
	err := runCargoLLVMCov("/nonexistent/cargo-binary", dir, out)
	if err == nil {
		t.Fatal("expected error for non-executable cargo binary")
	}
}

func TestRunCargoLLVMCovSuccess(t *testing.T) {
	binDir := t.TempDir()
	cargo := filepath.Join(binDir, "cargo")
	// A fake cargo that exits 0 covers the success path.
	if err := os.WriteFile(cargo, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "llvm-cov.json")
	if err := runCargoLLVMCov(cargo, dir, out); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestRunCargoLLVMCovNonZeroExit(t *testing.T) {
	binDir := t.TempDir()
	cargo := filepath.Join(binDir, "cargo")
	if err := os.WriteFile(cargo, []byte("#!/bin/sh\necho oops 1>&2\nexit 3\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "llvm-cov.json")
	err := runCargoLLVMCov(cargo, dir, out)
	if err == nil {
		t.Fatal("expected error when cargo exits non-zero")
	}
}
