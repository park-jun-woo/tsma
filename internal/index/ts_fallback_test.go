package index

import (
	"testing"
)

// TestNewTSTreeSitterIndexerFileFallback exercises the injected fileFallback
// closure (the line-based per-file re-index) directly.
func TestNewTSTreeSitterIndexerFileFallback(t *testing.T) {
	idx := newTSTreeSitterIndexer()
	funcs := idx.fileFallback("src/math.ts", "../../testdata/typescript/src/math.ts")
	if len(funcs) == 0 {
		t.Fatal("fileFallback returned no functions for math.ts")
	}
	found := false
	for _, f := range funcs {
		if f.Name == "add" {
			found = true
		}
	}
	if !found {
		t.Errorf("fileFallback did not find add: %+v", funcs)
	}
}

// TestTreeSitterIndexerUnavailableFallback forces tree-sitter to be unresolved
// (a bogus TSMA_TREE_SITTER name LookPath cannot find) so the constructor yields
// command=="" and Index delegates wholesale to the line-based fallback.
func TestTreeSitterIndexerUnavailableFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	idx := newTSTreeSitterIndexer()
	if idx.available() {
		t.Fatal("indexer should be unavailable with a bogus command")
	}
	funcs, err := idx.Index("../../testdata/typescript")
	if err != nil {
		t.Fatalf("Index (fallback): %v", err)
	}
	if len(funcs) == 0 {
		t.Fatal("line-based fallback produced no functions")
	}
}

// TestTreeSitterIndexerEmptyDir covers the len(files)==0 early return on the
// precise path (tree-sitter available, but the directory has no source files).
func TestTreeSitterIndexerEmptyDir(t *testing.T) {
	if _, _, ok := locateTreeSitter(t); !ok {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	idx := newTSTreeSitterIndexer()
	if !idx.available() {
		t.Fatal("indexer unavailable after locate")
	}
	funcs, err := idx.Index(t.TempDir())
	if err != nil {
		t.Fatalf("Index (empty dir): %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("empty dir produced %d funcs", len(funcs))
	}
}
