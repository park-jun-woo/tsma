package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// installFakeDotnet prepends a fake dotnet (with the given script body) to PATH.
func installFakeDotnet(t *testing.T, body string) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "dotnet"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// TestCsRunnerRunWithoutDotnet verifies a clear error when the .NET SDK is
// missing. A live dotnet run is not exercised here (toolchain required for E2E).
func TestCsRunnerRunWithoutDotnet(t *testing.T) {
	if _, err := exec.LookPath("dotnet"); err == nil {
		t.Skip("dotnet is installed; skipping toolchain-missing path")
	}

	r := &CsRunner{}
	_, err := r.Run("/tmp", "App.Tests/FooTests.cs")
	if err == nil {
		t.Fatal("expected error when dotnet is missing")
	}
	if !strings.Contains(err.Error(), "dotnet") {
		t.Errorf("error should mention dotnet: %v", err)
	}
}

func TestCsRunnerRunPass(t *testing.T) {
	installFakeDotnet(t, "exit 0\n")
	r := &CsRunner{}
	res, err := r.Run(t.TempDir(), "App.Tests/FooTests.cs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestCsRunnerRunFail(t *testing.T) {
	installFakeDotnet(t, "echo Failed! 1>&2\nexit 1\n")
	r := &CsRunner{}
	res, err := r.Run(t.TempDir(), "App.Tests/FooTests.cs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false when dotnet test fails")
	}
	if res.Output == "" {
		t.Error("expected non-empty output for failing run")
	}
}

func TestCsRunnerImplementsRunner(t *testing.T) {
	var _ Runner = &CsRunner{}
}
