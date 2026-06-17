package treesitter

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
)

// mathTS returns the absolute path to the shared TypeScript fixture.
func mathTS(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "..", "testdata", "typescript", "src", "math.ts"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	return p
}

func TestAttrValue(t *testing.T) {
	e := xml.StartElement{Attr: []xml.Attr{
		{Name: xml.Name{Local: "field"}, Value: "name"},
		{Name: xml.Name{Local: "srow"}, Value: "3"},
	}}
	if got := attrValue(e, "field"); got != "name" {
		t.Errorf("attrValue(field) = %q, want name", got)
	}
	if got := attrValue(e, "missing"); got != "" {
		t.Errorf("attrValue(missing) = %q, want empty", got)
	}
}

func TestAttrInt(t *testing.T) {
	e := xml.StartElement{Attr: []xml.Attr{
		{Name: xml.Name{Local: "srow"}, Value: "7"},
		{Name: xml.Name{Local: "scol"}, Value: "notanint"},
	}}
	if got := attrInt(e, "srow"); got != 7 {
		t.Errorf("attrInt(srow) = %d, want 7", got)
	}
	if got := attrInt(e, "scol"); got != 0 {
		t.Errorf("attrInt(invalid) = %d, want 0", got)
	}
	if got := attrInt(e, "missing"); got != 0 {
		t.Errorf("attrInt(missing) = %d, want 0", got)
	}
}

func TestChildByField(t *testing.T) {
	n := &Node{Children: []*Node{
		{Type: "identifier", Field: "name"},
		{Type: "statement_block"},
	}}
	if c := n.ChildByField("name"); c == nil || c.Type != "identifier" {
		t.Errorf("ChildByField(name) = %+v", c)
	}
	if c := n.ChildByField("nope"); c != nil {
		t.Errorf("ChildByField(nope) = %+v, want nil", c)
	}
}

func TestChildByType(t *testing.T) {
	n := &Node{Children: []*Node{
		{Type: "identifier", Field: "name"},
		{Type: "statement_block"},
	}}
	if c := n.ChildByType("statement_block"); c == nil {
		t.Error("ChildByType(statement_block) = nil")
	}
	if c := n.ChildByType("nope"); c != nil {
		t.Errorf("ChildByType(nope) = %+v, want nil", c)
	}
}

func TestStartEndLine(t *testing.T) {
	n := &Node{SRow: 4, ERow: 9}
	if got := n.StartLine(); got != 5 {
		t.Errorf("StartLine = %d, want 5", got)
	}
	if got := n.EndLine(); got != 10 {
		t.Errorf("EndLine = %d, want 10", got)
	}
}

func TestWalk(t *testing.T) {
	root := &Node{Type: "program", Children: []*Node{
		{Type: "a", Children: []*Node{{Type: "a1"}}},
		{Type: "b"},
	}}

	// nil node is a no-op.
	Walk(nil, func(*Node) bool { t.Fatal("visited nil"); return true })

	// full pre-order traversal.
	var seen []string
	Walk(root, func(n *Node) bool {
		seen = append(seen, n.Type)
		return true
	})
	want := []string{"program", "a", "a1", "b"}
	if len(seen) != len(want) {
		t.Fatalf("full walk = %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("full walk = %v, want %v", seen, want)
		}
	}

	// returning false prunes the subtree (a1 not visited).
	var pruned []string
	Walk(root, func(n *Node) bool {
		pruned = append(pruned, n.Type)
		return n.Type != "a"
	})
	for _, ty := range pruned {
		if ty == "a1" {
			t.Errorf("pruned walk visited a1: %v", pruned)
		}
	}
}

func TestGrammarDirExists(t *testing.T) {
	dir := t.TempDir()
	if !grammarDirExists(dir) {
		t.Error("grammarDirExists(tempdir) = false")
	}
	if grammarDirExists(filepath.Join(dir, "nope")) {
		t.Error("grammarDirExists(missing) = true")
	}
	f := filepath.Join(dir, "file")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if grammarDirExists(f) {
		t.Error("grammarDirExists(file) = true")
	}
}

func TestEnvGrammarDir(t *testing.T) {
	dir := t.TempDir()
	if got := envGrammarDir(dir); got != dir {
		t.Errorf("envGrammarDir(dir) = %q, want %q", got, dir)
	}
	if got := envGrammarDir(filepath.Join(dir, "nope")); got != "" {
		t.Errorf("envGrammarDir(missing) = %q, want empty", got)
	}
}

func TestGrammarCandidateBases(t *testing.T) {
	bases := grammarCandidateBases()
	if len(bases) < 2 || bases[0] != "." || bases[1] != "/tmp" {
		t.Fatalf("grammarCandidateBases = %v", bases)
	}
	// UserHomeDir normally succeeds in the test env, appending a third base.
	if home, err := os.UserHomeDir(); err == nil {
		if len(bases) != 3 || bases[2] != home {
			t.Errorf("expected home %q appended, got %v", home, bases)
		}
	}
}

func TestResolveCommand(t *testing.T) {
	// absolute, existing, non-dir override is returned verbatim.
	f := filepath.Join(t.TempDir(), "fake-ts")
	if err := os.WriteFile(f, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TSMA_TREE_SITTER", f)
	if got := ResolveCommand(); got != f {
		t.Errorf("abs override = %q, want %q", got, f)
	}

	// absolute, existing directory -> "".
	t.Setenv("TSMA_TREE_SITTER", t.TempDir())
	if got := ResolveCommand(); got != "" {
		t.Errorf("abs dir override = %q, want empty", got)
	}

	// absolute, missing path -> "".
	t.Setenv("TSMA_TREE_SITTER", "/no/such/abs/tree-sitter")
	if got := ResolveCommand(); got != "" {
		t.Errorf("abs missing override = %q, want empty", got)
	}

	// relative name not on PATH -> "".
	t.Setenv("TSMA_TREE_SITTER", "tsma-no-such-binary-xyz")
	if got := ResolveCommand(); got != "" {
		t.Errorf("missing PATH name = %q, want empty", got)
	}

	// default name resolved via PATH (requires tree-sitter on PATH).
	t.Setenv("TSMA_TREE_SITTER", "")
	if got := ResolveCommand(); got == "" {
		t.Skip("tree-sitter not on PATH; PATH-resolution branch not exercised")
	}
}

func TestResolveGrammar(t *testing.T) {
	// unknown language -> "".
	if got := ResolveGrammar("klingon"); got != "" {
		t.Errorf("ResolveGrammar(klingon) = %q, want empty", got)
	}

	// env override wins when it is an existing directory.
	dir := t.TempDir()
	t.Setenv("TSMA_TS_GRAMMAR", dir)
	if got := ResolveGrammar("typescript"); got != dir {
		t.Errorf("env override = %q, want %q", got, dir)
	}

	// bad env override falls through to "" (no probe match here).
	t.Setenv("TSMA_TS_GRAMMAR", filepath.Join(dir, "nope"))
	if got := ResolveGrammar("typescript"); got != "" {
		t.Errorf("bad env override = %q, want empty", got)
	}

	// no env override: probe node_modules bases. The test env installs the
	// grammar under /tmp, so the probe should resolve a directory.
	os.Unsetenv("TSMA_TS_GRAMMAR")
	probe := filepath.Join("/tmp", "node_modules", "tree-sitter-typescript", "typescript")
	if grammarDirExists(probe) {
		if got := ResolveGrammar("typescript"); got == "" {
			t.Errorf("probe expected to find %q, got empty", probe)
		}
	}
}

func TestRunAndParseFile(t *testing.T) {
	cmd := ResolveCommand()
	if cmd == "" {
		t.Skip("tree-sitter unavailable; CLI paths not exercised")
	}
	grammar := ResolveGrammar("typescript")
	path := mathTS(t)

	out, err := Run(cmd, grammar, []string{path})
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("Run produced no output")
	}

	root, err := ParseFile(cmd, grammar, path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if root == nil || root.Type != "program" {
		t.Fatalf("ParseFile root = %+v", root)
	}
}

func TestRunErrors(t *testing.T) {
	// bad command, empty output -> error returned.
	if _, err := Run("tsma-no-such-binary-xyz", "", []string{"x.ts"}); err == nil {
		t.Error("Run(bad command) expected error")
	}
}

func TestParseFileRunError(t *testing.T) {
	if _, err := ParseFile("tsma-no-such-binary-xyz", "", "x.ts"); err == nil {
		t.Error("ParseFile(bad command) expected error")
	}
}

// cannedBuilderXML exercises every branch of the xmlTreeBuilder handlers with no
// CLI dependency:
//   - an element before <sources> (handleStart !inSources early return; handleEnd
//     pop-skip with empty stack),
//   - nested nodes (attach to parent vs. as root),
//   - a node whose text precedes a child then resumes (handleChar "already set"),
//   - whitespace char data while the stack is empty (handleChar empty-stack return).
const cannedBuilderXML = `<?xml version="1.0"?>
<ignored>before</ignored>
<sources>
  <source name="/x/foo.ts">
    <program srow="0" scol="0" erow="3" ecol="0">
      <function_declaration srow="0" scol="0" erow="2" ecol="1">function<identifier field="name" srow="0" scol="9" erow="0" ecol="12">foo</identifier>tail</function_declaration>
    </program>
  </source>
</sources>`

func TestXMLBuilderBranches(t *testing.T) {
	sources, err := ParseXML([]byte(cannedBuilderXML))
	if err != nil {
		t.Fatalf("ParseXML error: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("want 1 source, got %d", len(sources))
	}
	src := sources[0]
	if src.Name != "/x/foo.ts" {
		t.Errorf("source name = %q", src.Name)
	}
	if src.Root == nil || src.Root.Type != "program" {
		t.Fatalf("root = %+v", src.Root)
	}
	fn := src.Root.ChildByType("function_declaration")
	if fn == nil {
		t.Fatal("no function_declaration")
	}
	// first non-empty run kept; "tail" after the child must be ignored.
	if fn.Text != "function" {
		t.Errorf("fn.Text = %q, want function", fn.Text)
	}
	name := fn.ChildByField("name")
	if name == nil || name.Text != "foo" {
		t.Errorf("name = %+v", name)
	}
}

func TestConsumeIgnoresUnknownTokens(t *testing.T) {
	// xml.ProcInst / comments are not start/end/chardata: consume must ignore them.
	b := &xmlTreeBuilder{}
	b.consume(xml.Comment("c"))
	b.consume(xml.ProcInst{Target: "xml", Inst: []byte("version=\"1.0\"")})
	if len(b.sources) != 0 || b.curRoot != nil {
		t.Errorf("consume mutated builder on unknown tokens: %+v", b)
	}
	// consume dispatches each token kind to its handler.
	b.consume(xml.StartElement{Name: xml.Name{Local: "sources"}})
	if !b.inSources {
		t.Error("consume(StartElement sources) did not set inSources")
	}
	b.consume(xml.StartElement{Name: xml.Name{Local: "source"}, Attr: []xml.Attr{{Name: xml.Name{Local: "name"}, Value: "/p.ts"}}})
	b.consume(xml.StartElement{Name: xml.Name{Local: "program"}})
	b.consume(xml.CharData("x"))
	b.consume(xml.EndElement{Name: xml.Name{Local: "program"}})
	b.consume(xml.EndElement{Name: xml.Name{Local: "source"}})
	b.consume(xml.EndElement{Name: xml.Name{Local: "sources"}})
	if len(b.sources) != 1 || b.sources[0].Name != "/p.ts" {
		t.Errorf("consume sequence built %+v", b.sources)
	}
}

// startEl is a small helper for building XML start elements in handler tests.
func startEl(name string, attrs ...string) xml.StartElement {
	e := xml.StartElement{Name: xml.Name{Local: name}}
	for i := 0; i+1 < len(attrs); i += 2 {
		e.Attr = append(e.Attr, xml.Attr{Name: xml.Name{Local: attrs[i]}, Value: attrs[i+1]})
	}
	return e
}

func TestHandleStart(t *testing.T) {
	b := &xmlTreeBuilder{}

	// <sources> sets the guard.
	b.handleStart(startEl("sources"))
	if !b.inSources {
		t.Fatal("inSources not set")
	}

	// a node before any <source> is still inSources but has empty stack -> root.
	b.handleStart(startEl("source", "name", "/a.ts"))
	if b.curSourceName != "/a.ts" || b.stack != nil || b.curRoot != nil {
		t.Fatalf("source reset failed: %+v", b)
	}

	// first real node attaches as the source root.
	b.handleStart(startEl("program", "srow", "0", "scol", "0", "erow", "5", "ecol", "0"))
	if b.curRoot == nil || b.curRoot.Type != "program" || len(b.stack) != 1 {
		t.Fatalf("program node not rooted: %+v", b.curRoot)
	}

	// nested node attaches under the stack top and pushes.
	b.handleStart(startEl("identifier", "field", "name", "srow", "1", "scol", "2", "erow", "1", "ecol", "5"))
	if len(b.stack) != 2 || len(b.curRoot.Children) != 1 {
		t.Fatalf("nested node not attached: %+v", b.curRoot)
	}
	child := b.curRoot.Children[0]
	if child.Field != "name" || child.SRow != 1 || child.SCol != 2 || child.ERow != 1 || child.ECol != 5 {
		t.Errorf("child fields = %+v", child)
	}

	// a node while !inSources is ignored.
	b2 := &xmlTreeBuilder{}
	b2.handleStart(startEl("stray"))
	if b2.curRoot != nil || len(b2.stack) != 0 {
		t.Errorf("stray node not ignored: %+v", b2)
	}
}

func TestHandleEnd(t *testing.T) {
	b := &xmlTreeBuilder{inSources: true, curSourceName: "/a.ts"}
	root := &Node{Type: "program"}
	b.curRoot = root
	b.stack = []*Node{root, {Type: "child"}}

	// generic end pops the stack.
	b.handleEnd(xml.EndElement{Name: xml.Name{Local: "child"}})
	if len(b.stack) != 1 {
		t.Fatalf("stack not popped: %d", len(b.stack))
	}

	// </source> finalizes the Source and resets.
	b.handleEnd(xml.EndElement{Name: xml.Name{Local: "source"}})
	if len(b.sources) != 1 || b.sources[0].Root != root || b.stack != nil || b.curRoot != nil {
		t.Fatalf("source not finalized: %+v", b)
	}

	// </sources> clears the guard.
	b.handleEnd(xml.EndElement{Name: xml.Name{Local: "sources"}})
	if b.inSources {
		t.Error("inSources not cleared")
	}

	// generic end with an empty stack is a no-op (no panic).
	b.handleEnd(xml.EndElement{Name: xml.Name{Local: "whatever"}})
}

func TestHandleChar(t *testing.T) {
	// empty stack -> ignored.
	b := &xmlTreeBuilder{}
	b.handleChar(xml.CharData("x"))

	n := &Node{Type: "identifier"}
	b.stack = []*Node{n}

	// whitespace-only is ignored.
	b.handleChar(xml.CharData("   \n  "))
	if n.Text != "" {
		t.Errorf("whitespace set text: %q", n.Text)
	}

	// first non-empty run is kept.
	b.handleChar(xml.CharData("  foo  "))
	if n.Text != "foo" {
		t.Errorf("text = %q, want foo", n.Text)
	}

	// a later run is ignored once text is set.
	b.handleChar(xml.CharData("bar"))
	if n.Text != "foo" {
		t.Errorf("text overwritten = %q, want foo", n.Text)
	}
}

func TestAttach(t *testing.T) {
	b := &xmlTreeBuilder{}

	// empty stack -> becomes the source root.
	root := &Node{Type: "program"}
	b.attach(root)
	if b.curRoot != root {
		t.Fatal("attach to empty stack did not set curRoot")
	}

	// non-empty stack -> appended under the top.
	b.stack = []*Node{root}
	child := &Node{Type: "identifier"}
	b.attach(child)
	if len(root.Children) != 1 || root.Children[0] != child {
		t.Fatalf("attach under top failed: %+v", root.Children)
	}
}
