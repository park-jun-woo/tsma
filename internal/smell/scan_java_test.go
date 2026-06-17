package smell

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

// locateJavaSmell finds a usable tree-sitter CLI + tree-sitter-java grammar and
// exports the env vars ScanJava reads. Returns false (caller t.Skip) when absent.
func locateJavaSmell(t *testing.T) bool {
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
	grammar := os.Getenv("TSMA_JAVA_GRAMMAR")
	if grammar == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", "tree-sitter-java")
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
	t.Setenv("TSMA_JAVA_GRAMMAR", grammar)
	return true
}

// TestScanJavaPositive asserts every escape hatch in the ReflectionTest fixture
// fires with the right rule (getDeclaredField/Method → JV-001, setAccessible(true)
// → JV-002) and nothing else does.
func TestScanJavaPositive(t *testing.T) {
	if !locateJavaSmell(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	findings, err := ScanJava("../../testdata/java/smell/ReflectionTest.java")
	if err != nil {
		t.Fatalf("ScanJava: %v", err)
	}
	counts := map[string]int{}
	for _, f := range findings {
		counts[f.Rule]++
	}
	// getDeclaredField + getDeclaredMethod
	if counts["TS-REFL-JV-001"] != 2 {
		t.Errorf("TS-REFL-JV-001 count = %d, want 2 (findings=%+v)", counts["TS-REFL-JV-001"], findings)
	}
	// two setAccessible(true) calls
	if counts["TS-REFL-JV-002"] != 2 {
		t.Errorf("TS-REFL-JV-002 count = %d, want 2 (findings=%+v)", counts["TS-REFL-JV-002"], findings)
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

// TestScanJavaNegative asserts the CleanTest fixture (public API + a string
// "getDeclaredMethod" argument + setAccessible(false)) produces zero findings —
// false-positive zero.
func TestScanJavaNegative(t *testing.T) {
	if !locateJavaSmell(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	findings, err := ScanJava("../../testdata/java/smell/CleanTest.java")
	if err != nil {
		t.Fatalf("ScanJava: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("clean fixture produced findings: %+v", findings)
	}
}
