package index

import "testing"

func TestNewIndexerReturnsGoIndexer(t *testing.T) {
	idx := NewIndexer("go")
	if _, ok := idx.(*GoIndexer); !ok {
		t.Errorf("NewIndexer(\"go\") returned %T, want *GoIndexer", idx)
	}
}

func TestNewIndexerReturnsTSIndexer(t *testing.T) {
	idx := NewIndexer("typescript")
	if _, ok := idx.(*TSIndexer); !ok {
		t.Errorf("NewIndexer(\"typescript\") returned %T, want *TSIndexer", idx)
	}
}

func TestNewIndexerReturnsPyIndexer(t *testing.T) {
	idx := NewIndexer("python")
	if _, ok := idx.(*PyIndexer); !ok {
		t.Errorf("NewIndexer(\"python\") returned %T, want *PyIndexer", idx)
	}
}

func TestNewIndexerReturnsUnsupported(t *testing.T) {
	idx := NewIndexer("rust")
	u, ok := idx.(*UnsupportedIndexer)
	if !ok {
		t.Fatalf("NewIndexer(\"rust\") returned %T, want *UnsupportedIndexer", idx)
	}
	if u.Lang != "rust" {
		t.Errorf("UnsupportedIndexer.Lang = %q, want %q", u.Lang, "rust")
	}
}
