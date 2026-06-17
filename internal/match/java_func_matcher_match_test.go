package match

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateJavaTS finds a usable tree-sitter CLI + tree-sitter-java grammar for the
// content-aware matcher test and exports the env vars. Returns false → t.Skip.
func locateJavaTS(t *testing.T) bool {
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

// TestJavaFuncMatcherContent asserts content-aware attribution across the JUnit
// src/main↔src/test mirror: a method called only by CalculatorTest (calc.add,
// calc.classify, new Calculator) attributes to CalculatorTest.java, and a
// StringUtils method attributes to StringUtilsTest.java — not by filename alone
// but by which test file actually invokes the symbol.
func TestJavaFuncMatcherContent(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	root := "../../testdata/java"
	m := &JavaFuncMatcher{}

	cases := []struct {
		fn   model.Function
		want string
	}{
		{model.Function{Name: "add", File: "src/main/java/com/example/calc/Calculator.java"},
			"src/test/java/com/example/calc/CalculatorTest.java"},
		{model.Function{Name: "classify", File: "src/main/java/com/example/calc/Calculator.java"},
			"src/test/java/com/example/calc/CalculatorTest.java"},
		{model.Function{Name: "repeat", File: "src/main/java/com/example/calc/StringUtils.java"},
			"src/test/java/com/example/calc/StringUtilsTest.java"},
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

// TestJavaFuncMatcherConstructorByNew asserts a constructor attributes via a
// `new Calculator(...)` object-creation in the test (content-based, not filename).
func TestJavaFuncMatcherConstructorByNew(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	m := &JavaFuncMatcher{}
	fn := model.Function{Name: "Calculator", File: "src/main/java/com/example/calc/Calculator.java"}
	tm, ok := m.MatchFunc("../../testdata/java", &fn)
	if !ok || len(tm.Files) == 0 {
		t.Fatalf("constructor not attributed: ok=%v tm=%+v", ok, tm)
	}
	if filepath.ToSlash(tm.Files[0]) != "src/test/java/com/example/calc/CalculatorTest.java" {
		t.Errorf("constructor attributed to %v", tm.Files)
	}
}
