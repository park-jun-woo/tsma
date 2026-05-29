package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRunCsCoverageBadBinary verifies the wrapper surfaces an error when the
// dotnet binary cannot be executed. A successful run requires a real .NET SDK +
// coverlet.collector (E2E only).
func TestRunCsCoverageBadBinary(t *testing.T) {
	err := runCsCoverage("/nonexistent/dotnet-binary", t.TempDir(), []string{"test"})
	if err == nil {
		t.Fatal("expected error for non-executable dotnet binary")
	}
}

func TestRunCsCoverageSuccess(t *testing.T) {
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "dotnet")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runCsCoverage(bin, t.TempDir(), []string{"test"}); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestRunCsCoverageNonZeroExit(t *testing.T) {
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "dotnet")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\necho oops 1>&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runCsCoverage(bin, t.TempDir(), []string{"test"}); err == nil {
		t.Fatal("expected error when dotnet exits non-zero")
	}
}
