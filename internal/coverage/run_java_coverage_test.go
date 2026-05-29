package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRunJavaCoverageBadBinary verifies the wrapper surfaces an error when the
// build-tool binary cannot be executed. A successful run requires a real JDK +
// build tool + JaCoCo plugin (E2E only).
func TestRunJavaCoverageBadBinary(t *testing.T) {
	err := runJavaCoverage("/nonexistent/mvn-binary", t.TempDir(), []string{"test"})
	if err == nil {
		t.Fatal("expected error for non-executable build-tool binary")
	}
}

func TestRunJavaCoverageSuccess(t *testing.T) {
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "mvn")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runJavaCoverage(bin, t.TempDir(), []string{"test", "jacoco:report"}); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestRunJavaCoverageNonZeroExit(t *testing.T) {
	binDir := t.TempDir()
	bin := filepath.Join(binDir, "mvn")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\necho oops 1>&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := runJavaCoverage(bin, t.TempDir(), []string{"test"}); err == nil {
		t.Fatal("expected error when build tool exits non-zero")
	}
}
