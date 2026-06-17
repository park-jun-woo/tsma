package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateCsTreeSitter finds a usable tree-sitter CLI + tree-sitter-c-sharp grammar
// for the integration test, probing env, PATH, and common node_modules bases. It
// returns ok=false (caller t.Skip) when either is absent — the precise path is an
// optional prerequisite (the plan's skip-gate). On success it exports
// TSMA_TREE_SITTER/TSMA_CSHARP_GRAMMAR for the test.
func locateCsTreeSitter(t *testing.T) bool {
	t.Helper()

	cmd := os.Getenv("TSMA_TREE_SITTER")
	if cmd == "" {
		if p, err := exec.LookPath("tree-sitter"); err == nil {
			cmd = p
		}
	}
	if cmd == "" {
		cmd = probe(tsLocateBases(), filepath.Join("node_modules", ".bin", "tree-sitter"))
	}

	grammar := os.Getenv("TSMA_CSHARP_GRAMMAR")
	if grammar == "" {
		grammar = probeDir(tsLocateBases(), filepath.Join("node_modules", "tree-sitter-c-sharp"))
	}

	if cmd == "" || grammar == "" {
		return false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_CSHARP_GRAMMAR", grammar)
	return true
}

// TestCSharpTreeSitterIndexerIntegration runs the real tree-sitter CLI over the
// testdata/csharp fixtures (skip-gated when CLI/grammar absent) and asserts
// precise method/constructor/property discovery — including the multi-line
// `Classify` signature the line-based regex cannot place, the nested Inner class,
// the file-scoped namespace qualification, and the attribute-spanning
// constructor. Test sources (*.Tests project) must never be indexed.
func TestCSharpTreeSitterIndexerIntegration(t *testing.T) {
	if !locateCsTreeSitter(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}

	idx := newCSharpTreeSitterIndexer()
	if !idx.available() {
		t.Fatal("indexer reports unavailable after locate")
	}

	funcs, err := idx.Index("../../testdata/csharp")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
		if strings.Contains(filepath.ToSlash(f.File), ".Tests/") {
			t.Errorf("test source should not be indexed: %+v", f)
		}
	}

	type want struct {
		start, end int
		exported   bool
	}
	wants := map[string]want{
		"Calc.Calculator.Total":      {7, 7, true},
		"Calc.Calculator.Calculator": {9, 13, true},
		"Calc.Calculator.Classify":   {15, 28, true},
		"Calc.Calculator.Helper":     {30, 33, false},
		"Calc.Calculator.Inner.Ping": {37, 39, true},
		"Calc.StringUtils.IsBlank":   {5, 12, true},
		"Calc.StringUtils.Repeat":    {14, 22, true},
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
}
