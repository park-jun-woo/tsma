package tsmagate

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// locateCsTSGate finds a usable tree-sitter CLI + tree-sitter-c-sharp grammar and
// exports the env vars smell.ScanCs reads. Returns false -> caller t.Skip.
func locateCsTSGate(t *testing.T) bool {
	t.Helper()
	bases := []string{".", "..", "../..", "/tmp"}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, home)
	}
	cmd := os.Getenv("TSMA_TREE_SITTER")
	if cmd == "" {
		if p, err := exec.LookPath("tree-sitter"); err == nil {
			cmd = p
		}
	}
	if cmd == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", ".bin", "tree-sitter")
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				cmd, _ = filepath.Abs(p)
				break
			}
		}
	}
	grammar := os.Getenv("TSMA_CSHARP_GRAMMAR")
	if grammar == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", "tree-sitter-c-sharp")
			if fi, err := os.Stat(p); err == nil && fi.IsDir() {
				grammar, _ = filepath.Abs(p)
				break
			}
		}
	}
	if cmd == "" || grammar == "" {
		return false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_CSHARP_GRAMMAR", grammar)
	return true
}

// TestScanCsSmellsContinue covers the error branch: when tree-sitter is
// unavailable every ScanCs errors and is skipped, yielding no findings.
func TestScanCsSmellsContinue(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	got := scanCsSmells("../../testdata/csharp", []string{"Calc.Tests/ReflectionTests.cs"})
	if len(got) != 0 {
		t.Errorf("expected no findings when tree-sitter unavailable, got %+v", got)
	}
}

// TestScanCsSmellsAppend covers the append branch with a real CLI: the
// ReflectionTests fixture yields escape-hatch findings, while a nonexistent file
// errors and is skipped (continue) within the same loop.
func TestScanCsSmellsAppend(t *testing.T) {
	if !locateCsTSGate(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	got := scanCsSmells("../../testdata/csharp", []string{
		"Calc.Tests/ReflectionTests.cs",
		"Calc.Tests/does_not_exist.cs", // errors -> continue
	})
	if len(got) == 0 {
		t.Fatal("expected findings from ReflectionTests.cs")
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-CS-001" && f.Rule != "TS-REFL-CS-002" {
			t.Errorf("unexpected rule %q", f.Rule)
		}
	}
}
