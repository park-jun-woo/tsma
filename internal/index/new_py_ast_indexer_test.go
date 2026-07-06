package index

import "testing"

// TestNewPyAstIndexerWiring covers the Python indexer factory: the line-based
// PyIndexer must be wired as the graceful fallback. The interpreter field is
// whatever resolvePython finds (possibly "" — Index then delegates to the
// fallback, which is exactly the wired contract).
func TestNewPyAstIndexerWiring(t *testing.T) {
	idx := newPyAstIndexer()
	if idx == nil {
		t.Fatal("newPyAstIndexer returned nil")
	}
	if idx.fallback == nil {
		t.Fatal("fallback must be wired")
	}
	if _, ok := idx.fallback.(*PyIndexer); !ok {
		t.Errorf("fallback = %T, want *PyIndexer", idx.fallback)
	}
}
