package match

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsSkipNoTreeSitter skips when the tree-sitter CLI/grammar is unavailable — the
// content-aware index needs them; the filename fallback is exercised separately.
func rsSkipNoTreeSitter(t *testing.T) {
	t.Helper()
	if treesitter.ResolveCommand() == "" || treesitter.ResolveGrammar("rust") == "" {
		t.Skip("tree-sitter CLI + rust grammar not available")
	}
}

// TestRsFuncMatcherInFile attributes a non-pub function (private_helper) to its
// own source file via the in-file #[cfg(test)] module — the case the tests/
// integration path cannot cover (external crate sees pub only).
func TestRsFuncMatcherInFile(t *testing.T) {
	rsSkipNoTreeSitter(t)
	root := "../../testdata/rust"
	fn := &model.Function{Name: "private_helper", File: "src/calc.rs"}
	tm, ok := (&RsFuncMatcher{}).MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected a content match for private_helper")
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "src/calc.rs" {
		t.Errorf("Files = %v, want [src/calc.rs]", tm.Files)
	}
}

// TestRsFuncMatcherIntegration attributes a pub function called only from the
// tests/ integration crate. `add` is also called in-file, so it resolves to the
// source file; `double` (nested::double) is referenced from integration.rs.
func TestRsFuncMatcherIntegration(t *testing.T) {
	rsSkipNoTreeSitter(t)
	root := "../../testdata/rust"
	fn := &model.Function{Name: "double", File: "src/calc.rs"}
	tm, ok := (&RsFuncMatcher{}).MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected a match for double")
	}
	var hasIntegration bool
	for _, f := range tm.Files {
		if filepath.ToSlash(f) == "tests/integration.rs" {
			hasIntegration = true
		}
	}
	if !hasIntegration {
		t.Errorf("Files = %v, want to include tests/integration.rs", tm.Files)
	}
}

// TestRsFuncMatcherFilenameFallback forces tree-sitter off so MatchFunc degrades
// to RsMatcher (in-file source file), proving the fallback is preserved.
func TestRsFuncMatcherFilenameFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	root := "../../testdata/rust"
	fn := &model.Function{Name: "add", File: "src/calc.rs"}
	tm, ok := (&RsFuncMatcher{}).MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected filename fallback to attribute add")
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "src/calc.rs" {
		t.Errorf("Files = %v, want [src/calc.rs] (in-file fallback)", tm.Files)
	}
}
