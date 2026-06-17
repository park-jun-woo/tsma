package smell

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

// locateCsSmell finds a usable tree-sitter CLI + tree-sitter-c-sharp grammar and
// exports the env vars ScanCs reads. Returns false (caller t.Skip) when absent.
func locateCsSmell(t *testing.T) bool {
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

// TestScanCsPositive asserts every escape hatch in the ReflectionTests fixture
// fires with the right rule (GetMethod/GetField/GetProperty → CS-001,
// MethodInfo/FieldInfo/PropertyInfo declarations → CS-002) and nothing else does.
func TestScanCsPositive(t *testing.T) {
	if !locateCsSmell(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	findings, err := ScanCs("../../testdata/csharp/Calc.Tests/ReflectionTests.cs")
	if err != nil {
		t.Fatalf("ScanCs: %v", err)
	}
	counts := map[string]int{}
	for _, f := range findings {
		counts[f.Rule]++
	}
	// GetMethod + GetField + GetProperty
	if counts["TS-REFL-CS-001"] != 3 {
		t.Errorf("TS-REFL-CS-001 count = %d, want 3 (findings=%+v)", counts["TS-REFL-CS-001"], findings)
	}
	// MethodInfo + FieldInfo + PropertyInfo declarations
	if counts["TS-REFL-CS-002"] != 3 {
		t.Errorf("TS-REFL-CS-002 count = %d, want 3 (findings=%+v)", counts["TS-REFL-CS-002"], findings)
	}
	var rules []string
	for r := range counts {
		rules = append(rules, r)
	}
	sort.Strings(rules)
	if len(rules) != 2 {
		t.Errorf("unexpected rules fired: %v", rules)
	}
}

// TestScanCsNegative asserts the CalculatorTests fixture (public API + a string
// "GetMethod" argument + a "MethodInfo" comment) produces zero findings —
// false-positive zero.
func TestScanCsNegative(t *testing.T) {
	if !locateCsSmell(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	findings, err := ScanCs("../../testdata/csharp/Calc.Tests/CalculatorTests.cs")
	if err != nil {
		t.Fatalf("ScanCs: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("clean fixture produced findings: %+v", findings)
	}
}
