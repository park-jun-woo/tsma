package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

func TestTSCalleeName(t *testing.T) {
	// bare identifier call.
	id := &treesitter.Node{Type: "identifier", Text: "foo"}
	if got := tsCalleeName(id); got != "foo" {
		t.Errorf("tsCalleeName(identifier) = %q, want foo", got)
	}
	// member call obj.bar -> property name.
	mem := &treesitter.Node{Type: "member_expression", Children: []*treesitter.Node{
		{Type: "identifier", Field: "object", Text: "obj"},
		{Type: "property_identifier", Field: "property", Text: "bar"},
	}}
	if got := tsCalleeName(mem); got != "bar" {
		t.Errorf("tsCalleeName(member) = %q, want bar", got)
	}
	// member call with no property field -> "".
	memNoProp := &treesitter.Node{Type: "member_expression"}
	if got := tsCalleeName(memNoProp); got != "" {
		t.Errorf("tsCalleeName(member w/o property) = %q, want empty", got)
	}
	// other node type -> "".
	if got := tsCalleeName(&treesitter.Node{Type: "call_expression"}); got != "" {
		t.Errorf("tsCalleeName(other) = %q, want empty", got)
	}
}

func TestCanonicalTSTestPathDirect(t *testing.T) {
	cases := []struct{ src, base, want string }{
		{"src/math.ts", "math.ts", "src/math.test.ts"},
		{"a/w.tsx", "w.tsx", "a/w.test.tsx"},
		{"lib/u.js", "u.js", "lib/u.test.js"},
		{"c.jsx", "c.jsx", "c.test.jsx"},
		{"src/x.go", "x.go", ""},
		{"r.md", "r.md", ""},
	}
	for _, c := range cases {
		if got := canonicalTSTestPath(c.src, c.base); got != c.want {
			t.Errorf("canonicalTSTestPath(%q,%q) = %q, want %q", c.src, c.base, got, c.want)
		}
	}
}

func TestTSFilenameFallback(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(rel string) {
		if err := os.WriteFile(filepath.Join(root, rel), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("src/widget.ts")
	write("src/widget.test.ts")
	write("src/lonely.ts")

	// a function with a sibling test file -> found.
	fn := &model.Function{Name: "widget", File: "src/widget.ts"}
	tm, ok := tsFilenameFallback(root, fn)
	if !ok {
		t.Fatal("tsFilenameFallback: want found for widget")
	}
	if len(tm.Files) != 1 || tm.Files[0] != filepath.Join("src", "widget.test.ts") {
		t.Errorf("Files = %v", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil (whole file)", tm.TestFuncs)
	}

	// a function without any sibling test -> not found.
	if _, ok := tsFilenameFallback(root, &model.Function{Name: "lonely", File: "src/lonely.ts"}); ok {
		t.Error("tsFilenameFallback: want not found for lonely")
	}
}

func TestMatchFuncNilFunc(t *testing.T) {
	m := &TypeScriptFuncMatcher{}
	if _, ok := m.MatchFunc("root", nil); ok {
		t.Error("MatchFunc(nil) = ok, want false")
	}
}

func TestMatchFuncFallbackNoTests(t *testing.T) {
	if !locateTS(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	// a package with a source file but no test file: the content index is nil, so
	// MatchFunc falls back to filename matching, which also finds nothing.
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "p"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "p", "foo.ts"), []byte("export function foo() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := &TypeScriptFuncMatcher{}
	if tm, ok := m.MatchFunc(root, &model.Function{Name: "foo", File: "p/foo.ts"}); ok {
		t.Errorf("MatchFunc(no tests) = %v, want not found", tm.Files)
	}

	// now add a sibling test file: filename fallback should attribute it even
	// though no test references foo by name.
	if err := os.WriteFile(filepath.Join(root, "p", "foo.test.ts"), []byte("it('x', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := m.MatchFunc(root, &model.Function{Name: "foo", File: "p/foo.ts"}); !ok {
		t.Error("MatchFunc: want filename fallback to find foo.test.ts")
	}
}

func TestIngestTSDirMissing(t *testing.T) {
	// a non-existent directory is a no-op (ReadDir errors -> return).
	idx := &TSPkgTestIndex{refs: map[string][]string{}}
	ingestTSDir(idx, t.TempDir(), "no-such-dir", "cmd", "grammar")
	if len(idx.refs) != 0 {
		t.Errorf("missing dir mutated idx: %v", idx.refs)
	}
}

func TestIngestTSTestFileParseError(t *testing.T) {
	// a bad command / missing file makes ParseFile error -> the file is skipped.
	idx := &TSPkgTestIndex{refs: map[string][]string{}}
	ingestTSTestFile(idx, "tsma-no-such-binary-xyz", "", "/no/such/file.ts", "file.ts")
	if len(idx.refs) != 0 {
		t.Errorf("parse error mutated idx: %v", idx.refs)
	}
}

func TestBuildTSPkgTestIndexUnavailable(t *testing.T) {
	// a bogus tree-sitter command -> ResolveCommand returns "" -> nil index.
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	if idx := BuildTSPkgTestIndex("../../testdata/typescript", "src"); idx != nil {
		t.Errorf("BuildTSPkgTestIndex(unavailable) = %+v, want nil", idx)
	}
}

func TestBuildTSPkgTestIndexContent(t *testing.T) {
	if !locateTS(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	// the src package's math.test.ts references add/classify.
	idx := BuildTSPkgTestIndex("../../testdata/typescript", "src")
	if idx == nil {
		t.Fatal("BuildTSPkgTestIndex returned nil for a package with tests")
	}
	for _, name := range []string{"add", "classify"} {
		if len(idx.refs[name]) == 0 {
			t.Errorf("no refs recorded for %q (have %v)", name, idx.refs)
		}
	}
}

func TestBuildTSPkgTestIndexNoTests(t *testing.T) {
	if !locateTS(t) {
		t.Skip("tree-sitter CLI + typescript grammar not available")
	}
	// a package directory with a source file but no *.test/*.spec file yields no
	// references, so the index is nil (MatchFunc then falls back to filenames).
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "p"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "p", "foo.ts"), []byte("export function foo() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if idx := BuildTSPkgTestIndex(root, "p"); idx != nil {
		t.Errorf("BuildTSPkgTestIndex(no tests) = %+v, want nil", idx)
	}
}
