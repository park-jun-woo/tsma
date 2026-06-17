package index

import "testing"

// TestNewRsTreeSitterIndexerFileFallback exercises the injected fileFallback
// closure (the line-based per-file re-index) directly, proving the legacy
// indexRsFile path is preserved under the precise indexer (zero regression).
func TestNewRsTreeSitterIndexerFileFallback(t *testing.T) {
	idx := newRsTreeSitterIndexer()
	funcs := idx.fileFallback("src/calc.rs", "../../testdata/rust/src/calc.rs")
	if len(funcs) == 0 {
		t.Fatal("fileFallback returned no functions for calc.rs")
	}
	for _, f := range funcs {
		if f.Name == "add_branches" || f.Name == "sub_works" {
			t.Errorf("line-based fallback indexed a #[cfg(test)] function: %+v", f)
		}
	}
}

// TestRsTreeSitterIndexerUnavailableFallback forces tree-sitter to be unresolved
// (a bogus TSMA_TREE_SITTER LookPath cannot find) so the constructor yields
// command=="" and Index delegates wholesale to the line-based RsIndexer — zero
// regression in environments without tree-sitter.
func TestRsTreeSitterIndexerUnavailableFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	idx := newRsTreeSitterIndexer()
	if idx.available() {
		t.Fatal("indexer should be unavailable with a bogus command")
	}
	funcs, err := idx.Index("../../testdata/rust")
	if err != nil {
		t.Fatalf("Index (fallback): %v", err)
	}
	if len(funcs) == 0 {
		t.Fatal("line-based fallback produced no functions")
	}
}
