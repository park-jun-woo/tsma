package match

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateCsTS finds a usable tree-sitter CLI + tree-sitter-c-sharp grammar for the
// content-aware matcher test and exports the env vars. Returns false → t.Skip.
func locateCsTS(t *testing.T) bool {
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

// TestCsFuncMatcherContent asserts content-aware attribution across the parallel
// *.Tests project: a method called only by CalculatorTests (calc.Classify, new
// Calculator) attributes to CalculatorTests.cs, and the StringUtils.Repeat method
// called by StringUtilsTests attributes to StringUtilsTests.cs — not by filename
// alone but by which test file actually invokes the symbol.
func TestCsFuncMatcherContent(t *testing.T) {
	if !locateCsTS(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	root := "../../testdata/csharp"
	m := &CsFuncMatcher{}

	cases := []struct {
		fn   model.Function
		want string
	}{
		{model.Function{Name: "Classify", File: "Calc/Calculator.cs"},
			"Calc.Tests/CalculatorTests.cs"},
		{model.Function{Name: "Repeat", File: "Calc/StringUtils.cs"},
			"Calc.Tests/StringUtilsTests.cs"},
	}

	for _, c := range cases {
		tm, ok := m.MatchFunc(root, &c.fn)
		if !ok {
			t.Errorf("%s: no match", c.fn.Name)
			continue
		}
		if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != c.want {
			t.Errorf("%s: files = %v, want [%s]", c.fn.Name, tm.Files, c.want)
		}
	}
}

// TestCsFuncMatcherConstructorByNew asserts a constructor attributes via a
// `new Calculator(...)` object-creation in the test (content-based, not filename).
func TestCsFuncMatcherConstructorByNew(t *testing.T) {
	if !locateCsTS(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	m := &CsFuncMatcher{}
	fn := model.Function{Name: "Calculator", File: "Calc/Calculator.cs"}
	tm, ok := m.MatchFunc("../../testdata/csharp", &fn)
	if !ok || len(tm.Files) == 0 {
		t.Fatalf("constructor not attributed: ok=%v tm=%+v", ok, tm)
	}
	if filepath.ToSlash(tm.Files[0]) != "Calc.Tests/CalculatorTests.cs" {
		t.Errorf("constructor attributed to %v", tm.Files)
	}
}
