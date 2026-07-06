package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
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

// stubIndexer is a canned whole-project fallback for hermetic Index tests.
type stubIndexer struct {
	funcs []model.Function
	err   error
}

func (s *stubIndexer) Index(string) ([]model.Function, error) { return s.funcs, s.err }

// writeFakeTreeSitter writes an executable shell script standing in for the
// tree-sitter CLI so every Index branch runs without the real binary.
func writeFakeTreeSitter(t *testing.T, script string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fake-tree-sitter")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

// newFakeIndexer builds a TreeSitterIndexer with injected command/extract/
// fallback closures for hermetic branch tests.
func newFakeIndexer(command string) (*TreeSitterIndexer, *stubIndexer) {
	stub := &stubIndexer{funcs: []model.Function{{Name: "fromProjectFallback"}}}
	return &TreeSitterIndexer{
		lang:     "fake",
		command:  command,
		fallback: stub,
		fileFallback: func(relPath, absPath string) []model.Function {
			return []model.Function{{Name: "fromFileFallback", File: relPath}}
		},
		isSource: func(p string) bool { return filepath.Ext(p) == ".ts" },
		skipDir:  func(string) error { return nil },
		extract: func(root *treesitter.Node, relPath, relDir string) []model.Function {
			return []model.Function{{Name: "fromExtract", File: relPath}}
		},
	}, stub
}

// TestTreeSitterIndexerUnavailableDelegates covers the available()==false
// guard: an empty command means the precise path cannot run, so Index
// delegates wholesale to the whole-project line-based fallback.
func TestTreeSitterIndexerUnavailableDelegates(t *testing.T) {
	idx, stub := newFakeIndexer("")
	funcs, err := idx.Index(t.TempDir())
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != stub.funcs[0].Name {
		t.Errorf("funcs = %+v, want project fallback", funcs)
	}
}

// TestTreeSitterIndexerCollectError forces collectSourceFiles to fail (skipDir
// returns a real error) so Index delegates to the whole-project fallback.
func TestTreeSitterIndexerCollectError(t *testing.T) {
	idx, stub := newFakeIndexer("fake-command-never-run")
	idx.skipDir = func(string) error { return os.ErrPermission }
	funcs, err := idx.Index(t.TempDir())
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != stub.funcs[0].Name {
		t.Errorf("funcs = %+v, want project fallback", funcs)
	}
}

// TestTreeSitterIndexerNoFiles covers the len(files)==0 early return without
// requiring the real CLI (the command is never invoked for an empty project).
func TestTreeSitterIndexerNoFiles(t *testing.T) {
	idx, _ := newFakeIndexer("fake-command-never-run")
	funcs, err := idx.Index(t.TempDir())
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if funcs != nil {
		t.Errorf("funcs = %+v, want nil", funcs)
	}
}

// TestTreeSitterIndexerBatchFallback drives the two whole-batch fallback
// branches: a CLI failure with no output (Run error) and syntactically broken
// XML (ParseXML error). Both must re-index every file line-based.
func TestTreeSitterIndexerBatchFallback(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{"run error empty output", "exit 1\n"},
		{"malformed xml", "printf '<sources><open></close></sources>'\nexit 0\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "a.ts"), []byte("export function a() {}\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			idx, _ := newFakeIndexer(writeFakeTreeSitter(t, tt.script))
			funcs, err := idx.Index(root)
			if err != nil {
				t.Fatalf("Index: %v", err)
			}
			if len(funcs) != 1 || funcs[0].Name != "fromFileFallback" {
				t.Errorf("funcs = %+v, want one fromFileFallback", funcs)
			}
		})
	}
}

// TestTreeSitterIndexerExtractAndPerFileFallback feeds a batch where the fake
// CLI parsed a.ts but omitted b.ts: a.ts must go through extract and b.ts
// through the per-file line-based fallback.
func TestTreeSitterIndexerExtractAndPerFileFallback(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.ts", "b.ts"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("export function f() {}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	xml := "<sources>\n<source name=\"" + filepath.Join(root, "a.ts") + "\">\n" +
		"<program srow=\"0\" scol=\"0\" erow=\"1\" ecol=\"0\"></program>\n</source>\n</sources>\n"
	script := "cat <<'TSMA_EOF'\n" + xml + "TSMA_EOF\nexit 0\n"
	idx, _ := newFakeIndexer(writeFakeTreeSitter(t, script))

	funcs, err := idx.Index(root)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	got := map[string]string{}
	for _, f := range funcs {
		got[f.File] = f.Name
	}
	want := map[string]string{"a.ts": "fromExtract", "b.ts": "fromFileFallback"}
	if len(funcs) != 2 {
		t.Fatalf("funcs = %+v, want 2", funcs)
	}
	for file, name := range want {
		if got[file] != name {
			t.Errorf("%s indexed by %q, want %q (all: %v)", file, got[file], name, got)
		}
	}
}
