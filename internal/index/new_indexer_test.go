package index

import "testing"

func TestNewIndexerReturnsGoIndexer(t *testing.T) {
	idx := NewIndexer("go")
	if _, ok := idx.(*GoIndexer); !ok {
		t.Errorf("NewIndexer(\"go\") returned %T, want *GoIndexer", idx)
	}
}

func TestNewIndexerReturnsTSIndexer(t *testing.T) {
	// Phase005a: the typescript factory entry returns the tree-sitter precise
	// indexer (TSIndexer is retained as its fallback, asserted in index_test.go).
	idx := NewIndexer("typescript")
	if _, ok := idx.(*TreeSitterIndexer); !ok {
		t.Errorf("NewIndexer(\"typescript\") returned %T, want *TreeSitterIndexer", idx)
	}
}

func TestNewIndexerReturnsPyIndexer(t *testing.T) {
	// Phase005b: the python factory entry returns the ast precise indexer
	// (PyIndexer is retained as its line-based fallback, asserted in index_test.go).
	idx := NewIndexer("python")
	if _, ok := idx.(*PyAstIndexer); !ok {
		t.Errorf("NewIndexer(\"python\") returned %T, want *PyAstIndexer", idx)
	}
}

func TestNewIndexerReturnsRsIndexer(t *testing.T) {
	idx := NewIndexer("rust")
	if _, ok := idx.(*RsIndexer); !ok {
		t.Errorf("NewIndexer(\"rust\") returned %T, want *RsIndexer", idx)
	}
}

func TestNewIndexerReturnsJavaIndexer(t *testing.T) {
	// Phase005c: the java factory entry returns the tree-sitter precise indexer
	// (JavaIndexer is retained as its line-based fallback, asserted in
	// java_fallback_test.go).
	idx := NewIndexer("java")
	if _, ok := idx.(*TreeSitterIndexer); !ok {
		t.Errorf("NewIndexer(\"java\") returned %T, want *TreeSitterIndexer", idx)
	}
}

func TestNewIndexerReturnsCsIndexer(t *testing.T) {
	// Phase005d: the csharp factory entry returns the tree-sitter precise indexer
	// (CsIndexer is retained as its line-based fallback, asserted in
	// cs_fallback_test.go).
	idx := NewIndexer("csharp")
	if _, ok := idx.(*TreeSitterIndexer); !ok {
		t.Errorf("NewIndexer(\"csharp\") returned %T, want *TreeSitterIndexer", idx)
	}
}

func TestNewIndexerReturnsUnsupported(t *testing.T) {
	idx := NewIndexer("kotlin")
	u, ok := idx.(*UnsupportedIndexer)
	if !ok {
		t.Fatalf("NewIndexer(\"kotlin\") returned %T, want *UnsupportedIndexer", idx)
	}
	if u.Lang != "kotlin" {
		t.Errorf("UnsupportedIndexer.Lang = %q, want %q", u.Lang, "kotlin")
	}
}
