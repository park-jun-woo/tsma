package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// locateRsTreeSitter finds a usable tree-sitter CLI + tree-sitter-rust grammar
// for the integration test, probing env, PATH, and common node_modules bases. It
// returns ok=false (caller t.Skip) when either is absent — the precise path is an
// optional prerequisite (the plan's skip-gate).
func locateRsTreeSitter(t *testing.T) bool {
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

	grammar := os.Getenv("TSMA_RUST_GRAMMAR")
	if grammar == "" {
		grammar = probeDir(tsLocateBases(), filepath.Join("node_modules", "tree-sitter-rust"))
	}

	if cmd == "" || grammar == "" {
		return false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_RUST_GRAMMAR", grammar)
	return true
}

// TestRsTreeSitterIndexerIntegration runs the real tree-sitter CLI over the
// testdata/rust fixtures (skip-gated when CLI/grammar absent) and asserts precise
// fn/impl-method/generic/nested-module discovery — including the multi-line
// `compute` signature the line-based regex cannot place, the non-pub
// private_helper (a first-class target), and the nested `nested::double`
// qualification. The in-file #[cfg(test)] mod functions must NEVER be indexed.
func TestRsTreeSitterIndexerIntegration(t *testing.T) {
	if !locateRsTreeSitter(t) {
		t.Skip("tree-sitter CLI + rust grammar not available")
	}

	idx := newRsTreeSitterIndexer()
	if !idx.available() {
		t.Fatal("indexer reports unavailable after locate")
	}

	funcs, err := idx.Index("../../testdata/rust")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
		switch f.Name {
		case "add_branches", "sub_works", "private_helper_is_reachable",
			"compute_folds", "max_and_double", "unsafe_block_cheese",
			"transmute_cheese", "ptr_cheese":
			t.Errorf("#[cfg(test)] function should not be indexed: %+v", f)
		}
	}

	type want struct {
		start, end int
		exported   bool
	}
	wants := map[string]want{
		"src.add":                {6, 12, true},
		"src.sub":                {15, 17, true},
		"src.private_helper":     {21, 27, false},
		"src.max_of":             {30, 39, true},
		"src.Calculator.new":     {48, 50, true},
		"src.Calculator.compute": {55, 63, true},
		"src.nested.double":      {69, 71, true},
		"src.to_bytes":           {8, 10, true}, // cheese.rs, safe production fn
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
