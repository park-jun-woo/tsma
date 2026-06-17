package match

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestJavaFuncMatcherNilFn covers the nil-function guard.
func TestJavaFuncMatcherNilFn(t *testing.T) {
	if _, ok := (&JavaFuncMatcher{}).MatchFunc("../../testdata/java", nil); ok {
		t.Error("nil fn should report not found")
	}
}

// TestJavaFuncMatcherIndexNilFallback forces tree-sitter absent so the content
// index is nil and MatchFunc must fall back to the conventional FooTest.java.
func TestJavaFuncMatcherIndexNilFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	fn := &model.Function{Name: "add", File: "src/main/java/com/example/calc/Calculator.java"}
	tm, ok := (&JavaFuncMatcher{}).MatchFunc("../../testdata/java", fn)
	if !ok {
		t.Fatal("expected filename fallback when index is nil")
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "src/test/java/com/example/calc/CalculatorTest.java" {
		t.Errorf("fallback Files = %v", tm.Files)
	}
}

// TestJavaFuncMatcherContentMissFallback covers the CLI-present path where the
// index exists but does not reference the function, so it still falls back to
// the filename convention.
func TestJavaFuncMatcherContentMissFallback(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	fn := &model.Function{Name: "totallyAbsentMethod", File: "src/main/java/com/example/calc/Calculator.java"}
	tm, ok := (&JavaFuncMatcher{}).MatchFunc("../../testdata/java", fn)
	if !ok {
		t.Fatal("expected filename fallback on content miss")
	}
	if filepath.ToSlash(tm.Files[0]) != "src/test/java/com/example/calc/CalculatorTest.java" {
		t.Errorf("fallback Files = %v", tm.Files)
	}
}
