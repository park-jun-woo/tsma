package smell

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

func locateTSSmell(t *testing.T) bool {
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
	grammar := os.Getenv("TSMA_TS_GRAMMAR")
	if grammar == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", "tree-sitter-typescript", "typescript")
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
	t.Setenv("TSMA_TS_GRAMMAR", grammar)
	return true
}

// TestScanTSPositive asserts every escape hatch in the cheats fixture fires with
// the right rule, and nothing else does.
func TestScanTSPositive(t *testing.T) {
	if !locateTSSmell(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	findings, err := ScanTS("../../testdata/typescript/smell/cheats.test.ts")
	if err != nil {
		t.Fatalf("ScanTS: %v", err)
	}
	counts := map[string]int{}
	for _, f := range findings {
		counts[f.Rule]++
	}
	want := map[string]int{
		"TS-REFL-TS-001": 2,
		"TS-REFL-TS-002": 2,
		"TS-REFL-TS-003": 2,
	}
	var rules []string
	for r := range counts {
		rules = append(rules, r)
	}
	sort.Strings(rules)
	for rule, n := range want {
		if counts[rule] != n {
			t.Errorf("rule %s fired %d times, want %d (all: %v)", rule, counts[rule], n, rules)
		}
	}
	if len(counts) != len(want) {
		t.Errorf("unexpected rules fired: %v", rules)
	}
}

// TestScanTSNegative asserts the clean fixture produces zero findings.
func TestScanTSNegative(t *testing.T) {
	if !locateTSSmell(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	findings, err := ScanTS("../../testdata/typescript/smell/clean.test.ts")
	if err != nil {
		t.Fatalf("ScanTS: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("clean fixture must yield 0 findings, got %v", findings)
	}
}
