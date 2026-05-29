package coverage

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeCargo installs a fake `cargo` executable on PATH whose script body is the
// supplied shell snippet. It returns the project dir to run Check against. The
// fake bin dir is prepended to PATH (so the script can still use coreutils).
func fakeCargo(t *testing.T, script string) string {
	t.Helper()
	binDir := t.TempDir()
	cargoPath := filepath.Join(binDir, "cargo")
	if err := os.WriteFile(cargoPath, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return t.TempDir()
}

// TestRsCheckerCheckWithoutCargo verifies a clear error when the cargo
// toolchain is missing. A live coverage run requires cargo + cargo-llvm-cov
// (E2E only; not exercised in sandbox environments).
func TestRsCheckerCheckWithoutCargo(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err == nil {
		t.Skip("cargo is installed; skipping toolchain-missing path")
	}

	c := &RsChecker{}
	_, err := c.Check("/tmp", mkMatch("src/lib.rs"), &model.Function{File: "src/lib.rs", StartLine: 1, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when cargo is missing")
	}
	if !strings.Contains(err.Error(), "cargo") {
		t.Errorf("error should mention cargo: %v", err)
	}
}

func TestRsCheckerCheck_cargoRunFails(t *testing.T) {
	// Fake cargo exits non-zero: findCargo + MkdirAll succeed, then
	// runCargoLLVMCov returns an error -> "rust coverage failed" branch.
	proj := fakeCargo(t, "exit 1\n")

	c := &RsChecker{}
	_, err := c.Check(proj, mkMatch("src/lib.rs"), &model.Function{File: "src/lib.rs", StartLine: 1, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when cargo llvm-cov fails")
	}
	if !strings.Contains(err.Error(), "rust coverage failed") {
		t.Errorf("expected 'rust coverage failed', got: %v", err)
	}
}

func TestRsCheckerCheck_parseFails(t *testing.T) {
	// Fake cargo succeeds but writes invalid JSON to the output path so the
	// parse step fails.
	proj := fakeCargo(t, `
out=""
while [ $# -gt 0 ]; do
  if [ "$1" = "--output-path" ]; then out="$2"; fi
  shift
done
echo "{not json" > "$out"
exit 0
`)

	c := &RsChecker{}
	_, err := c.Check(proj, mkMatch("src/lib.rs"), &model.Function{File: "src/lib.rs", StartLine: 1, EndLine: 5})
	if err == nil {
		t.Fatal("expected parse error for invalid llvm-cov json")
	}
	if !strings.Contains(err.Error(), "parse llvm-cov json") {
		t.Errorf("expected 'parse llvm-cov json', got: %v", err)
	}
}

func TestRsCheckerCheck_success(t *testing.T) {
	// Locate the real fixture before changing PATH/dirs.
	fixtureAbs, err := filepath.Abs(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	// Fake cargo writes the valid fixture JSON to the requested output path.
	proj := fakeCargo(t, `
out=""
while [ $# -gt 0 ]; do
  if [ "$1" = "--output-path" ]; then out="$2"; fi
  shift
done
cat "`+fixtureAbs+`" > "$out"
exit 0
`)

	c := &RsChecker{}
	report, err := c.Check(proj, mkMatch("src/lib.rs"), &model.Function{File: "src/lib.rs", StartLine: 1, EndLine: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
}

func TestRsCheckerCheck_mkdirFails(t *testing.T) {
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "cargo"), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)

	// Make the project root's .tsma a file so MkdirAll fails.
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := &RsChecker{}
	_, err := c.Check(proj, mkMatch("src/lib.rs"), &model.Function{File: "src/lib.rs", StartLine: 1, EndLine: 5})
	if err == nil {
		t.Fatal("expected error when .tsma dir cannot be created")
	}
}

func TestRsCheckerImplementsChecker(t *testing.T) {
	var _ Checker = &RsChecker{}
}
