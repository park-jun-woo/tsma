package index

import "testing"

// TestNewJavaTreeSitterIndexerFileFallback exercises the injected fileFallback
// closure (the line-based per-file re-index) directly, proving the legacy
// indexJavaFile path is preserved under the precise indexer.
func TestNewJavaTreeSitterIndexerFileFallback(t *testing.T) {
	idx := newJavaTreeSitterIndexer()
	funcs := idx.fileFallback(
		"src/main/java/com/example/calc/Calculator.java",
		"../../testdata/java/src/main/java/com/example/calc/Calculator.java")
	if len(funcs) == 0 {
		t.Fatal("fileFallback returned no functions for Calculator.java")
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

// TestJavaTreeSitterIndexerUnavailableFallback forces tree-sitter to be
// unresolved (a bogus TSMA_TREE_SITTER LookPath cannot find) so the constructor
// yields command=="" and Index delegates wholesale to the line-based JavaIndexer
// — zero regression in environments without tree-sitter.
func TestJavaTreeSitterIndexerUnavailableFallback(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	idx := newJavaTreeSitterIndexer()
	if idx.available() {
		t.Fatal("indexer should be unavailable with a bogus command")
	}
	funcs, err := idx.Index("../../testdata/java")
	if err != nil {
		t.Fatalf("Index (fallback): %v", err)
	}
	if len(funcs) == 0 {
		t.Fatal("line-based fallback produced no functions")
	}
}
