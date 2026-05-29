package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// installFakeCargo prepends a fake cargo (with the given script body) to PATH.
func installFakeCargo(t *testing.T, body string) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "cargo"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// TestRsRunnerRunWithoutCargo verifies a clear error when the cargo toolchain
// is missing. A live cargo run is not exercised here (toolchain required for E2E).
func TestRsRunnerRunWithoutCargo(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err == nil {
		t.Skip("cargo is installed; skipping toolchain-missing path")
	}

	r := &RsRunner{}
	_, err := r.Run("/tmp", "tests/api.rs")
	if err == nil {
		t.Fatal("expected error when cargo is missing")
	}
	if !strings.Contains(err.Error(), "cargo") {
		t.Errorf("error should mention cargo: %v", err)
	}
}

func TestRsRunnerRun_pass(t *testing.T) {
	installFakeCargo(t, "exit 0\n")
	r := &RsRunner{}
	res, err := r.Run(t.TempDir(), "tests/api.rs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestRsRunnerRun_fail(t *testing.T) {
	installFakeCargo(t, "echo failures 1>&2\nexit 101\n")
	r := &RsRunner{}
	res, err := r.Run(t.TempDir(), "tests/api.rs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false when cargo test fails")
	}
	if res.Output == "" {
		t.Error("expected non-empty output for failing run")
	}
}

func TestRsRunnerImplementsRunner(t *testing.T) {
	var _ Runner = &RsRunner{}
}
