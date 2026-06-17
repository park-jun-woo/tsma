package smell

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsLocate finds a usable tree-sitter CLI + rust grammar for the smell tests, or
// skips (the precise path is an optional prerequisite — the plan's skip-gate).
func rsLocate(t *testing.T) (string, string) {
	t.Helper()
	cmd := treesitter.ResolveCommand()
	grammar := treesitter.ResolveGrammar("rust")
	if cmd == "" || grammar == "" {
		t.Skip("tree-sitter CLI + rust grammar not available")
	}
	return cmd, grammar
}

// TestScanRsCheesePositives asserts the three Rust escape-hatch rules fire on the
// in-file #[cfg(test)] module of cheese.rs (unsafe, transmute, ptr) — node-based,
// only inside the test scope.
func TestScanRsCheesePositives(t *testing.T) {
	rsLocate(t)
	findings, err := ScanRs(filepath.FromSlash("../../testdata/rust/src/cheese.rs"))
	if err != nil {
		t.Fatalf("ScanRs: %v", err)
	}
	got := map[string]int{}
	for _, f := range findings {
		got[f.Rule]++
	}
	for _, rule := range []string{"TS-REFL-RS-001", "TS-REFL-RS-002", "TS-REFL-RS-003"} {
		if got[rule] == 0 {
			t.Errorf("expected %s to fire on cheese.rs, findings=%+v", rule, findings)
		}
	}
}

// TestScanRsCalcCleanNoFindings asserts calc.rs — whose in-file #[cfg(test)]
// module is clean — yields zero findings (false-positive zero). Its production
// `to_bytes`-style safe code never trips a detector because scanning is scoped to
// the test module.
func TestScanRsCalcCleanNoFindings(t *testing.T) {
	rsLocate(t)
	findings, err := ScanRs(filepath.FromSlash("../../testdata/rust/src/calc.rs"))
	if err != nil {
		t.Fatalf("ScanRs: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings on clean calc.rs, got %+v", findings)
	}
}

// TestScanRsIntegrationClean asserts the integration test file (whole-file test
// scope) is also clean.
func TestScanRsIntegrationClean(t *testing.T) {
	rsLocate(t)
	findings, err := ScanRs(filepath.FromSlash("../../testdata/rust/tests/integration.rs"))
	if err != nil {
		t.Fatalf("ScanRs: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings on clean integration.rs, got %+v", findings)
	}
}
