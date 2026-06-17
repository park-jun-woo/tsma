package match

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsCfgAttrNode builds a #[cfg(test)] attribute_item subtree.
func rsCfgAttrNode() *treesitter.Node {
	return &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
		{Type: "identifier", Text: "cfg"},
		{Type: "token_tree", Children: []*treesitter.Node{
			{Type: "identifier", Text: "test"},
		}},
	}}
}

// TestCanonicalRsTestPath covers the .rs (source-as-test) and non-.rs branches.
func TestCanonicalRsTestPath(t *testing.T) {
	if got := canonicalRsTestPath("src/lib.rs", "lib.rs"); got != "src/lib.rs" {
		t.Errorf(".rs: canonicalRsTestPath = %q, want src/lib.rs", got)
	}
	if got := canonicalRsTestPath("src/lib.go", "lib.go"); got != "" {
		t.Errorf("non-.rs: canonicalRsTestPath = %q, want empty", got)
	}
}

// TestRsAttrIsCfgTest covers cfg+test (true), cfg-only, test-only, and neither.
func TestRsAttrIsCfgTest(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want bool
	}{
		{"cfg(test)", rsCfgAttrNode(), true},
		{"cfg only", &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "cfg"}}}, false},
		{"test only", &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "test"}}}, false},
		{"neither", &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "derive"}}}, false},
	}
	for _, c := range cases {
		if got := rsAttrIsCfgTest(c.node); got != c.want {
			t.Errorf("%s: rsAttrIsCfgTest = %v, want %v", c.desc, got, c.want)
		}
	}
}

// TestRsInvokedName covers nil, identifier, scoped_identifier (with/without name
// field), field_expression (with/without field), and the default branch.
func TestRsInvokedName(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want string
	}{
		{"nil", nil, ""},
		{"identifier", &treesitter.Node{Type: "identifier", Text: "free_fn"}, "free_fn"},
		{
			desc: "scoped_identifier with name",
			node: &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
				{Type: "identifier", Text: "module"},
				{Type: "identifier", Field: "name", Text: "new"},
			}},
			want: "new",
		},
		{"scoped_identifier without name", &treesitter.Node{Type: "scoped_identifier"}, ""},
		{
			desc: "field_expression with field",
			node: &treesitter.Node{Type: "field_expression", Children: []*treesitter.Node{
				{Type: "identifier", Field: "value", Text: "obj"},
				{Type: "field_identifier", Field: "field", Text: "method"},
			}},
			want: "method",
		},
		{"field_expression without field", &treesitter.Node{Type: "field_expression"}, ""},
		{"default", &treesitter.Node{Type: "closure_expression"}, ""},
	}
	for _, c := range cases {
		if got := rsInvokedName(c.node); got != c.want {
			t.Errorf("%s: rsInvokedName = %q, want %q", c.desc, got, c.want)
		}
	}
}

// TestRsMacroCallNames covers an identifier+token_tree pair (collected), a
// trailing identifier with no following token_tree (skipped), and a token_tree
// too short to iterate.
func TestRsMacroCallNames(t *testing.T) {
	tt := &treesitter.Node{Type: "token_tree", Children: []*treesitter.Node{
		{Type: "identifier", Text: "assert_eq"}, // followed by nested token_tree
		{Type: "token_tree", Children: []*treesitter.Node{
			{Type: "identifier", Text: "add"},
			{Type: "token_tree"},
		}},
		{Type: "identifier", Text: "trailing"}, // no following sibling
	}}
	got := rsMacroCallNames(tt)
	if !reflect.DeepEqual(got, []string{"assert_eq"}) {
		t.Errorf("rsMacroCallNames = %v, want [assert_eq]", got)
	}

	// token_tree with a single child: loop never runs.
	short := &treesitter.Node{Type: "token_tree", Children: []*treesitter.Node{{Type: "identifier", Text: "x"}}}
	if got := rsMacroCallNames(short); got != nil {
		t.Errorf("short token_tree: rsMacroCallNames = %v, want nil", got)
	}
}

// TestCollectRsCalledNames covers call_expression (named + empty-name skip),
// the macro token_tree branch, and an unrelated node type.
func TestCollectRsCalledNames(t *testing.T) {
	root := &treesitter.Node{Type: "block", Children: []*treesitter.Node{
		// bare call: free_fn()
		{Type: "call_expression", Children: []*treesitter.Node{
			{Type: "identifier", Field: "function", Text: "free_fn"},
		}},
		// call whose function resolves to "" (closure) -> not added.
		{Type: "call_expression", Children: []*treesitter.Node{
			{Type: "closure_expression", Field: "function"},
		}},
		// macro: assert_eq!(add(...))
		{Type: "macro_invocation", Children: []*treesitter.Node{
			{Type: "token_tree", Children: []*treesitter.Node{
				{Type: "identifier", Text: "add"},
				{Type: "token_tree"},
			}},
		}},
		// unrelated.
		{Type: "let_declaration"},
	}}
	got := collectRsCalledNames(root)
	want := map[string]struct{}{"free_fn": {}, "add": {}}
	if !reflect.DeepEqual(got, want) {
		gotKeys := keysOf(got)
		t.Errorf("collectRsCalledNames keys = %v, want [add free_fn]", gotKeys)
	}
}

func keysOf(m map[string]struct{}) []string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// TestRsCollectCfgTestBody covers the attribute (cfg(test) and non-test),
// mod_item (pending+body appended, pending+no-body, no-pending), and default.
func TestRsCollectCfgTestBody(t *testing.T) {
	body := &treesitter.Node{Type: "declaration_list", Field: "body"}
	modWithBody := &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{body}}
	modNoBody := &treesitter.Node{Type: "mod_item"}

	// cfg(test) attribute sets pending.
	var out []*treesitter.Node
	if p := rsCollectCfgTestBody(rsCfgAttrNode(), false, &out); !p {
		t.Error("cfg(test) attr: pending = false, want true")
	}
	// non-test attribute keeps pending unchanged.
	if p := rsCollectCfgTestBody(&treesitter.Node{Type: "attribute_item"}, false, &out); p {
		t.Error("non-test attr: pending = true, want false")
	}
	// mod_item with pending + body: appended, returns false.
	out = nil
	if p := rsCollectCfgTestBody(modWithBody, true, &out); p {
		t.Error("guarded mod: pending = true, want false")
	}
	if len(out) != 1 || out[0] != body {
		t.Errorf("guarded mod: out = %v, want [body]", out)
	}
	// mod_item with pending but no body: not appended.
	out = nil
	rsCollectCfgTestBody(modNoBody, true, &out)
	if len(out) != 0 {
		t.Errorf("body-less mod: out = %v, want empty", out)
	}
	// mod_item without pending: not appended.
	out = nil
	rsCollectCfgTestBody(modWithBody, false, &out)
	if len(out) != 0 {
		t.Errorf("unguarded mod: out = %v, want empty", out)
	}
	// default node carries pending.
	if p := rsCollectCfgTestBody(&treesitter.Node{Type: "use_declaration"}, true, &out); !p {
		t.Error("default: pending = false, want true (carried)")
	}
}

// TestRsCfgTestBodies covers nil root, a guarded mod (body collected), an
// unguarded mod (skipped), and an empty root.
func TestRsCfgTestBodies(t *testing.T) {
	if got := rsCfgTestBodies(nil); got != nil {
		t.Errorf("nil root: rsCfgTestBodies = %v, want nil", got)
	}

	body := &treesitter.Node{Type: "declaration_list", Field: "body"}
	root := &treesitter.Node{Type: "source_file", Children: []*treesitter.Node{
		rsCfgAttrNode(),
		{Type: "mod_item", Children: []*treesitter.Node{body}},
		// unguarded mod (no preceding cfg(test)): not collected.
		{Type: "mod_item", Children: []*treesitter.Node{
			{Type: "declaration_list", Field: "body"},
		}},
	}}
	got := rsCfgTestBodies(root)
	if len(got) != 1 || got[0] != body {
		t.Errorf("rsCfgTestBodies = %v, want [body]", got)
	}

	if got := rsCfgTestBodies(&treesitter.Node{Type: "source_file"}); got != nil {
		t.Errorf("empty root: rsCfgTestBodies = %v, want nil", got)
	}
}

// TestRsRefsToTestMatch covers the empty (found=false), dedup (order-preserving),
// and single-file branches.
func TestRsRefsToTestMatch(t *testing.T) {
	if _, ok := rsRefsToTestMatch(nil); ok {
		t.Error("empty: found = true, want false")
	}
	tm, ok := rsRefsToTestMatch([]string{"a.rs", "b.rs", "a.rs", "b.rs"})
	if !ok {
		t.Fatal("non-empty: found = false, want true")
	}
	if !reflect.DeepEqual(tm.Files, []string{"a.rs", "b.rs"}) {
		t.Errorf("Files = %v, want [a.rs b.rs] (deduped, order-preserving)", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil", tm.TestFuncs)
	}
}

// TestRsFilenameFallback covers the matched (in-file source) and unmatched (no
// test) branches, no tree-sitter needed.
func TestRsFilenameFallback(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("src", "lib.rs")
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("pub fn f() {}\n#[cfg(test)]\nmod tests {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fn := &model.Function{Name: "f", File: rel}
	tm, ok := rsFilenameFallback(dir, fn)
	if !ok || len(tm.Files) != 1 || tm.Files[0] != rel {
		t.Errorf("matched: rsFilenameFallback = %v, ok=%v, want [%s]", tm.Files, ok, rel)
	}

	// no test file present -> found=false.
	dir2 := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir2, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, rel), []byte("pub fn f() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := rsFilenameFallback(dir2, fn); ok {
		t.Error("no test: rsFilenameFallback found = true, want false")
	}
}

// TestMatchFuncNilAndContentMiss covers the nil-fn guard and the content-miss
// (index built but fn name absent) fallback. The nil branch needs no tree-sitter.
func TestMatchFuncNilAndContentMiss(t *testing.T) {
	if _, ok := (&RsFuncMatcher{}).MatchFunc("root", nil); ok {
		t.Error("nil fn: MatchFunc ok = true, want false")
	}

	rsSkipNoTreeSitter(t)
	// calc.rs has an in-file #[cfg(test)] module, so the index is non-nil, but a
	// name it never calls falls through to the filename fallback (the source file).
	fn := &model.Function{Name: "nonexistent_fn", File: "src/calc.rs"}
	tm, ok := (&RsFuncMatcher{}).MatchFunc("../../testdata/rust", fn)
	if !ok {
		t.Fatal("content miss: expected filename fallback to attribute the source file")
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "src/calc.rs" {
		t.Errorf("content miss: Files = %v, want [src/calc.rs]", tm.Files)
	}
}

// TestBuildRsTestIndexNoTreeSitter covers the command=="" nil branch (forced via
// a bogus absolute TSMA_TREE_SITTER).
func TestBuildRsTestIndexNoTreeSitter(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/tsma-bogus-tree-sitter")
	if idx := BuildRsTestIndex("../../testdata/rust", "src/calc.rs"); idx != nil {
		t.Errorf("no tree-sitter: BuildRsTestIndex = %v, want nil", idx)
	}
}

// TestBuildRsTestIndexReal covers the full build (in-file + tests/ refs) and the
// len(refs)==0 nil branch over a tree-less temp source file.
func TestBuildRsTestIndexReal(t *testing.T) {
	rsSkipNoTreeSitter(t)

	idx := BuildRsTestIndex("../../testdata/rust", "src/calc.rs")
	if idx == nil {
		t.Fatal("real build: BuildRsTestIndex = nil, want non-nil index")
	}
	// add is called both in-file and from tests/integration.rs.
	if len(idx.refs["add"]) == 0 {
		t.Errorf("refs[add] = %v, want non-empty", idx.refs["add"])
	}
	// double is called from tests/integration.rs (nested::double).
	var fromIntegration bool
	for _, f := range idx.refs["double"] {
		if filepath.ToSlash(f) == "tests/integration.rs" {
			fromIntegration = true
		}
	}
	if !fromIntegration {
		t.Errorf("refs[double] = %v, want to include tests/integration.rs", idx.refs["double"])
	}

	// A source file with no in-file tests and no tests/ refs -> len(refs)==0 -> nil.
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("src", "lonely.rs")
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("pub fn lonely() -> i32 { 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := BuildRsTestIndex(dir, rel); got != nil {
		t.Errorf("no refs: BuildRsTestIndex = %v, want nil", got)
	}
}

// TestIngestRsParseFailures covers the parse-failure no-op branches of
// ingestRsInFile and ingestRsTestFile (bogus command -> ParseFile error).
func TestIngestRsParseFailures(t *testing.T) {
	idx := &RsTestIndex{refs: make(map[string][]string)}
	ingestRsInFile(idx, "/nonexistent/bogus-ts", "", "/no/such/file.rs", "src/x.rs")
	ingestRsTestFile(idx, "/nonexistent/bogus-ts", "", "/no/such/file.rs", "tests/x.rs")
	if len(idx.refs) != 0 {
		t.Errorf("parse failures: refs = %v, want empty", idx.refs)
	}
}

// TestIngestRsTestsDir covers the missing-dir no-op, the IsDir / non-.rs skip
// branches, and the .rs ingest call. The directory traversal runs without
// tree-sitter (the per-file parse just no-ops), then the real-grammar path is
// asserted skip-gated.
func TestIngestRsTestsDir(t *testing.T) {
	// missing tests/ dir: no-op.
	idx := &RsTestIndex{refs: make(map[string][]string)}
	ingestRsTestsDir(idx, t.TempDir(), "/nonexistent/bogus-ts", "")
	if len(idx.refs) != 0 {
		t.Errorf("missing dir: refs = %v, want empty", idx.refs)
	}

	// dir with a subdir, a non-.rs file, and a .rs file: only the .rs reaches
	// ingestRsTestFile (which no-ops without tree-sitter), exercising the filters.
	root := t.TempDir()
	testsDir := filepath.Join(root, "tests")
	if err := os.MkdirAll(filepath.Join(testsDir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "it.rs"), []byte("// test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx2 := &RsTestIndex{refs: make(map[string][]string)}
	ingestRsTestsDir(idx2, root, "/nonexistent/bogus-ts", "")
	if len(idx2.refs) != 0 {
		t.Errorf("filtered dir (no tree-sitter): refs = %v, want empty", idx2.refs)
	}

	// real grammar: ingest the fixture tests/ dir and assert a ref appears.
	rsSkipNoTreeSitter(t)
	idx3 := &RsTestIndex{refs: make(map[string][]string)}
	ingestRsTestsDir(idx3, "../../testdata/rust", treesitter.ResolveCommand(), treesitter.ResolveGrammar("rust"))
	if len(idx3.refs["double"]) == 0 {
		t.Errorf("real tests dir: refs[double] = %v, want non-empty", idx3.refs["double"])
	}
}

// TestIngestRsInFileReal covers the success path of ingestRsInFile directly.
func TestIngestRsInFileReal(t *testing.T) {
	rsSkipNoTreeSitter(t)
	idx := &RsTestIndex{refs: make(map[string][]string)}
	abs, err := filepath.Abs("../../testdata/rust/src/calc.rs")
	if err != nil {
		t.Fatal(err)
	}
	ingestRsInFile(idx, treesitter.ResolveCommand(), treesitter.ResolveGrammar("rust"), abs, "src/calc.rs")
	if len(idx.refs["add"]) == 0 {
		t.Errorf("ingestRsInFile: refs[add] = %v, want non-empty", idx.refs["add"])
	}
}
