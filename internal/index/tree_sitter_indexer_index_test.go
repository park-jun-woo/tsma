package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestTreeSitterIndexerIntegration runs the real tree-sitter CLI over the
// testdata/typescript fixtures (skip-gated when the CLI/grammar is absent) and
// asserts precise function discovery — including the multi-line `scale`
// signature the line-based regex cannot match.
func TestTreeSitterIndexerIntegration(t *testing.T) {
	if _, _, ok := locateTreeSitter(t); !ok {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}

	idx := newTSTreeSitterIndexer()
	if !idx.available() {
		t.Fatal("indexer reports unavailable after locate")
	}

	funcs, err := idx.Index("../../testdata/typescript")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
		if f.File == "src/math.test.ts" {
			t.Errorf("test file should not be indexed: %+v", f)
		}
	}

	type want struct {
		start, end int
		exported   bool
	}
	wants := map[string]want{
		"src.add":             {3, 8, true},
		"src.classify":        {10, 18, true},
		"src.internalHelper":  {20, 22, false},
		"src.double":          {24, 26, true},
		"src.Rectangle.Area":  {13, 15, true},
		"src.Rectangle.scale": {17, 21, false},
		"src.makeSquare":      {24, 26, true},
		"src.perimeter":       {28, 30, true},
	}

	for qn, w := range wants {
		f, ok := byQN[qn]
		if !ok {
			t.Errorf("missing function %q (have %d funcs)", qn, len(funcs))
			continue
		}
		if f.StartLine != w.start || f.EndLine != w.end {
			t.Errorf("%s: range = %d..%d, want %d..%d", qn, f.StartLine, f.EndLine, w.start, w.end)
		}
		if f.Exported != w.exported {
			t.Errorf("%s: Exported = %v, want %v", qn, f.Exported, w.exported)
		}
	}

	// The constructor must never be indexed as a testable function.
	if _, ok := byQN["src.Rectangle.constructor"]; ok {
		t.Error("constructor should be skipped")
	}
}
