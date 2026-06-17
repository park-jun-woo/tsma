package match

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateTS finds the tree-sitter CLI + TS grammar for integration tests, setting
// the env the resolvers read; returns false to skip when unavailable.
func locateTS(t *testing.T) bool {
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

// TestTypeScriptFuncMatcherContent proves content-aware attribution: math.test.ts
// calls add/classify so those attribute to it, while makeSquare (called by no
// test, with no shapes.test.ts) is correctly NOT attributed — a naive
// "any test in the dir" matcher would wrongly claim makeSquare is tested.
func TestTypeScriptFuncMatcherContent(t *testing.T) {
	if !locateTS(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	root := "../../testdata/typescript"
	m := &TypeScriptFuncMatcher{}

	for _, name := range []string{"add", "classify"} {
		fn := &model.Function{Name: name, File: "src/math.ts"}
		tm, ok := m.MatchFunc(root, fn)
		if !ok {
			t.Errorf("%s: want content match, got none", name)
			continue
		}
		if len(tm.Files) != 1 || tm.Files[0] != "src/math.test.ts" {
			t.Errorf("%s: Files = %v, want [src/math.test.ts]", name, tm.Files)
		}
	}

	// makeSquare is referenced by no test and has no sibling test file.
	fn := &model.Function{Name: "makeSquare", File: "src/shapes.ts"}
	if tm, ok := m.MatchFunc(root, fn); ok {
		t.Errorf("makeSquare: want no attribution (content precision), got %v", tm.Files)
	}
}
