package match

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestJavaFilenameFallback covers both the found (conventional FooTest.java
// exists in the src/test mirror) and not-found branches against the testdata
// Maven layout — no tree-sitter CLI needed (pure os.Stat filename matching).
func TestJavaFilenameFallback(t *testing.T) {
	root := "../../testdata/java"

	// Found: Calculator.java -> src/test mirror has CalculatorTest.java.
	fn := &model.Function{Name: "add", File: "src/main/java/com/example/calc/Calculator.java"}
	tm, ok := javaFilenameFallback(root, fn)
	if !ok {
		t.Fatalf("expected fallback match for %s", fn.File)
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "src/test/java/com/example/calc/CalculatorTest.java" {
		t.Errorf("Files = %v", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil (run whole class)", tm.TestFuncs)
	}

	// Not found: no test mirror file for this source.
	missing := &model.Function{Name: "x", File: "src/main/java/com/example/calc/Nonexistent.java"}
	if _, ok := javaFilenameFallback(root, missing); ok {
		t.Error("expected no fallback for a source without a test mirror")
	}
}
