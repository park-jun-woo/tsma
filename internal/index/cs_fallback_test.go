package index

import "testing"

// TestNewCSharpTreeSitterIndexerFileFallback exercises the injected fileFallback
// closure (the line-based per-file re-index) directly, proving the legacy
// indexCsFile path is preserved under the precise indexer.
func TestNewCSharpTreeSitterIndexerFileFallback(t *testing.T) {
	idx := newCSharpTreeSitterIndexer()
	funcs := idx.fileFallback(
		"Calc/Calculator.cs",
		"../../testdata/csharp/Calc/Calculator.cs")
	if len(funcs) == 0 {
		t.Fatal("fileFallback returned no functions for Calculator.cs")
	}
	// Helper is a single-line signature the line-based path resolves; the
	// multi-line Classify is exactly what only the tree-sitter path catches.
	found := false
	for _, f := range funcs {
		if f.Name == "Helper" {
			found = true
		}
	}
	if !found {
		t.Errorf("fileFallback did not find Helper: %+v", funcs)
	}
}

// TestCSharpTreeSitterIndexerUnavailableFallback forces tree-sitter to be
// unresolved (a bogus TSMA_TREE_SITTER LookPath cannot find) so the constructor
// yields command=="" and Index delegates wholesale to the line-based CsIndexer
// — zero regression in environments without tree-sitter.
func TestCSharpTreeSitterIndexerUnavailableFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	idx := newCSharpTreeSitterIndexer()
	if idx.available() {
		t.Fatal("indexer should be unavailable with a bogus command")
	}
	funcs, err := idx.Index("../../testdata/csharp")
	if err != nil {
		t.Fatalf("Index (fallback): %v", err)
	}
	if len(funcs) == 0 {
		t.Fatal("line-based fallback produced no functions")
	}
}
