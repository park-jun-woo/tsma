package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateJavaTreeSitter finds a usable tree-sitter CLI + tree-sitter-java grammar
// for the integration test, probing env, PATH, and common node_modules bases. It
// returns ok=false (caller t.Skip) when either is absent — the precise path is an
// optional prerequisite (the plan's skip-gate). On success it exports
// TSMA_TREE_SITTER/TSMA_JAVA_GRAMMAR for the test.
func locateJavaTreeSitter(t *testing.T) bool {
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

	grammar := os.Getenv("TSMA_JAVA_GRAMMAR")
	if grammar == "" {
		grammar = probeDir(tsLocateBases(), filepath.Join("node_modules", "tree-sitter-java"))
	}

	if cmd == "" || grammar == "" {
		return false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_JAVA_GRAMMAR", grammar)
	return true
}

// TestJavaTreeSitterIndexerIntegration runs the real tree-sitter CLI over the
// testdata/java fixtures (skip-gated when CLI/grammar absent) and asserts precise
// method/constructor discovery — including the multi-line generic `classify`
// signature the line-based regex cannot place, the nested Helper class, and the
// package-qualified names. Test sources (src/test) must never be indexed.
func TestJavaTreeSitterIndexerIntegration(t *testing.T) {
	if !locateJavaTreeSitter(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}

	idx := newJavaTreeSitterIndexer()
	if !idx.available() {
		t.Fatal("indexer reports unavailable after locate")
	}

	funcs, err := idx.Index("../../testdata/java")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
		if strings.Contains(filepath.ToSlash(f.File), "src/test/") {
			t.Errorf("test source should not be indexed: %+v", f)
		}
	}

	type want struct {
		start, end int
		exported   bool
	}
	wants := map[string]want{
		"com.example.calc.Calculator.Calculator":       {7, 9, true},
		"com.example.calc.Calculator.add":               {11, 13, true},
		"com.example.calc.Calculator.classify":          {15, 24, true},
		"com.example.calc.Calculator.Helper.helperOnly": {27, 29, false},
		"com.example.calc.StringUtils.StringUtils":      {5, 6, false},
		"com.example.calc.StringUtils.isBlank":          {8, 10, true},
		"com.example.calc.StringUtils.repeat":           {12, 18, true},
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
