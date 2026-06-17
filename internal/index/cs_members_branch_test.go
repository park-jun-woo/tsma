package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// csNameField builds a `name` field identifier leaf node.
func csNameField(text string) *treesitter.Node {
	return &treesitter.Node{Type: "identifier", Field: "name", Text: text}
}

// csModifier builds a `(modifier)` leaf with the given keyword.
func csModifier(kw string) *treesitter.Node {
	return &treesitter.Node{Type: "modifier", Text: kw}
}

// csMethodNode builds a method_declaration with optional modifier children and a
// name field, spanning rows sRow..eRow (0-based).
func csMethodNode(name string, mods []string, sRow, eRow int) *treesitter.Node {
	n := &treesitter.Node{Type: "method_declaration", SRow: sRow, ERow: eRow}
	for _, m := range mods {
		n.Children = append(n.Children, csModifier(m))
	}
	if name != "" {
		n.Children = append(n.Children, csNameField(name))
	}
	return n
}

// TestCsNodeExported covers the public, non-public, and no-modifier branches,
// plus the type-mismatch guard (a "public"-text node that is not a modifier).
func TestCsNodeExported(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want bool
	}{
		{
			desc: "public modifier among several",
			node: &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{
				csModifier("public"), csModifier("static"),
			}},
			want: true,
		},
		{
			desc: "modifiers without public",
			node: &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{
				csModifier("private"), csModifier("static"),
			}},
			want: false,
		},
		{
			desc: "no modifier children",
			node: &treesitter.Node{Type: "method_declaration"},
			want: false,
		},
		{
			desc: "public text on a non-modifier node is ignored",
			node: &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{
				{Type: "identifier", Text: "public"},
			}},
			want: false,
		},
	}
	for _, c := range cases {
		if got := csNodeExported(c.node); got != c.want {
			t.Errorf("%s: csNodeExported = %v, want %v", c.desc, got, c.want)
		}
	}
}

// TestCsDottedName covers the nil, leaf-text, and qualified_name (joined
// identifier leaves) branches.
func TestCsDottedName(t *testing.T) {
	if got := csDottedName(nil); got != "" {
		t.Errorf("nil node: csDottedName = %q, want empty", got)
	}

	leaf := &treesitter.Node{Type: "identifier", Text: "Foo"}
	if got := csDottedName(leaf); got != "Foo" {
		t.Errorf("leaf: csDottedName = %q, want Foo", got)
	}

	// qualified_name with nested identifier leaves and one stray empty-text
	// identifier (Text != "" false) that must be skipped.
	qual := &treesitter.Node{Type: "qualified_name", Children: []*treesitter.Node{
		{Type: "identifier", Text: "A"},
		{Type: "identifier", Text: "B"},
		{Type: "identifier", Text: ""}, // skipped
		{Type: "."}, // non-identifier, skipped
		{Type: "identifier", Text: "C"},
	}}
	if got := csDottedName(qual); got != "A.B.C" {
		t.Errorf("qualified: csDottedName = %q, want A.B.C", got)
	}
}

// TestCsFileNamespace covers the no-file-scoped-namespace branch and the
// present file-scoped namespace (dotted) branch.
func TestCsFileNamespace(t *testing.T) {
	// no file_scoped_namespace_declaration child -> "".
	root := &treesitter.Node{Type: "compilation_unit", Children: []*treesitter.Node{
		{Type: "class_declaration"},
	}}
	if got := csFileNamespace(root); got != "" {
		t.Errorf("no namespace: csFileNamespace = %q, want empty", got)
	}

	// file-scoped namespace A.B (qualified_name name field).
	ns := &treesitter.Node{Type: "file_scoped_namespace_declaration", Children: []*treesitter.Node{
		{Type: "qualified_name", Field: "name", Children: []*treesitter.Node{
			{Type: "identifier", Text: "A"},
			{Type: "identifier", Text: "B"},
		}},
	}}
	root2 := &treesitter.Node{Type: "compilation_unit", Children: []*treesitter.Node{ns}}
	if got := csFileNamespace(root2); got != "A.B" {
		t.Errorf("file-scoped: csFileNamespace = %q, want A.B", got)
	}
}

// TestCsFuncFromNode covers the nameless / empty-name false branches and the
// successful conversion (span, qualified name, exported).
func TestCsFuncFromNode(t *testing.T) {
	scopes := []csScope{{typeName: "Outer"}}

	// nameless: no `name` field at all (e.g. operator declaration).
	if _, ok := csFuncFromNode(&treesitter.Node{Type: "method_declaration"}, "Ns", scopes, "F.cs"); ok {
		t.Error("nameless declaration should return ok=false")
	}

	// empty name: `name` field present but blank text.
	emptyName := &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{csNameField("")}}
	if _, ok := csFuncFromNode(emptyName, "Ns", scopes, "F.cs"); ok {
		t.Error("empty-name declaration should return ok=false")
	}

	// valid public method.
	node := csMethodNode("Compute", []string{"public"}, 4, 6)
	fn, ok := csFuncFromNode(node, "Ns", scopes, "F.cs")
	if !ok {
		t.Fatal("valid declaration should return ok=true")
	}
	if fn.QualifiedName != "Ns.Outer.Compute" || fn.Name != "Compute" ||
		fn.File != "F.cs" || fn.StartLine != 5 || fn.EndLine != 7 ||
		!fn.Exported || fn.Status != model.StatusTodo {
		t.Errorf("csFuncFromNode = %+v", fn)
	}
}

// TestWalkCSTypeDecl covers the nameless, empty-name, body-less and full
// (recurse) branches.
func TestWalkCSTypeDecl(t *testing.T) {
	// nameless type: skipped.
	var out []model.Function
	walkCSTypeDecl(&treesitter.Node{Type: "class_declaration"}, "Ns", nil, "F.cs", &out)
	if len(out) != 0 {
		t.Errorf("nameless type: out = %+v, want empty", out)
	}

	// named type with no body: returns after building scope, no funcs.
	out = nil
	noBody := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{csNameField("NoBody")}}
	walkCSTypeDecl(noBody, "Ns", nil, "F.cs", &out)
	if len(out) != 0 {
		t.Errorf("body-less type: out = %+v, want empty", out)
	}

	// named type with a body containing one method: recurse and qualify with the
	// pre-existing scope stack preserved.
	out = nil
	body := &treesitter.Node{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
		csMethodNode("M", []string{"public"}, 1, 1),
	}}
	withBody := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{
		csNameField("Widget"), body,
	}}
	walkCSTypeDecl(withBody, "Ns", []csScope{{typeName: "Top"}}, "F.cs", &out)
	if len(out) != 1 || out[0].QualifiedName != "Ns.Top.Widget.M" {
		t.Errorf("named-with-body: out = %+v, want single Ns.Top.Widget.M", out)
	}
}

// TestDispatchCSMember covers every recognized member shape, the recursive
// declaration shapes, the ignored file-scoped namespace, and the default case.
func TestDispatchCSMember(t *testing.T) {
	methodChild := func() *treesitter.Node { return csMethodNode("M", []string{"public"}, 1, 1) }
	typeDecl := func(nodeType, name string) *treesitter.Node {
		return &treesitter.Node{Type: nodeType, Children: []*treesitter.Node{
			csNameField(name),
			{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}
	}

	cases := []struct {
		desc string
		node *treesitter.Node
		want int // number of funcs emitted
	}{
		{"method_declaration", csMethodNode("Am", nil, 1, 1), 1},
		{"constructor_declaration", &treesitter.Node{Type: "constructor_declaration", Children: []*treesitter.Node{csNameField("Ctor")}}, 1},
		{"destructor_declaration", &treesitter.Node{Type: "destructor_declaration", Children: []*treesitter.Node{csNameField("Dtor")}}, 1},
		{"property_declaration", &treesitter.Node{Type: "property_declaration", Children: []*treesitter.Node{csNameField("Prop")}}, 1},
		{"namespace_declaration recurses", typeDecl("namespace_declaration", "Ns"), 1},
		{"class_declaration recurses", typeDecl("class_declaration", "C"), 1},
		{"struct_declaration recurses", typeDecl("struct_declaration", "S"), 1},
		{"interface_declaration recurses", typeDecl("interface_declaration", "I"), 1},
		{"record_declaration recurses", typeDecl("record_declaration", "R"), 1},
		{"record_struct_declaration recurses", typeDecl("record_struct_declaration", "RS"), 1},
		{"enum_declaration recurses", typeDecl("enum_declaration", "E"), 1},
		{"file_scoped_namespace_declaration ignored", &treesitter.Node{Type: "file_scoped_namespace_declaration"}, 0},
		{"field_declaration ignored", &treesitter.Node{Type: "field_declaration"}, 0},
	}
	for _, c := range cases {
		var out []model.Function
		dispatchCSMember(c.node, "", nil, "F.cs", &out)
		if len(out) != c.want {
			t.Errorf("%s: emitted %d funcs (%+v), want %d", c.desc, len(out), out, c.want)
		}
	}
}

// TestCollectCSMembers covers iteration over a container's children, including
// the empty-container (loop body never runs) branch.
func TestCollectCSMembers(t *testing.T) {
	container := &treesitter.Node{Type: "declaration_list", Children: []*treesitter.Node{
		csMethodNode("A", nil, 1, 1),
		csMethodNode("B", nil, 2, 2),
		{Type: "field_declaration"}, // ignored
	}}
	var out []model.Function
	collectCSMembers(container, "", nil, "F.cs", &out)
	if len(out) != 2 {
		t.Fatalf("collectCSMembers emitted %d funcs (%+v), want 2", len(out), out)
	}

	// empty container: loop body never runs.
	out = nil
	collectCSMembers(&treesitter.Node{Type: "declaration_list"}, "", nil, "F.cs", &out)
	if len(out) != 0 {
		t.Errorf("empty container: out = %+v, want empty", out)
	}
}

// TestExtractCSharpFunctions covers the nil-root guard and the full extraction
// path (file-scoped namespace + nested type qualification) over a hand-built
// compilation_unit, exercising the C# grammar interpreter without a subprocess.
func TestExtractCSharpFunctions(t *testing.T) {
	// nil root -> nil.
	if got := extractCSharpFunctions(nil, "F.cs", "."); got != nil {
		t.Errorf("nil root: extractCSharpFunctions = %+v, want nil", got)
	}

	// file-scoped namespace App, class Calc with a method and a nested class.
	ns := &treesitter.Node{Type: "file_scoped_namespace_declaration", Children: []*treesitter.Node{
		csNameField("App"),
	}}
	inner := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{
		csNameField("Inner"),
		{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			csMethodNode("Ping", []string{"public"}, 5, 5),
		}},
	}}
	outer := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{
		csNameField("Calc"),
		{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			csMethodNode("Total", []string{"public"}, 2, 2),
			inner,
		}},
	}}
	root := &treesitter.Node{Type: "compilation_unit", Children: []*treesitter.Node{ns, outer}}

	funcs := extractCSharpFunctions(root, "F.cs", ".")
	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
	}
	if _, ok := byQN["App.Calc.Total"]; !ok {
		t.Errorf("missing App.Calc.Total in %+v", funcs)
	}
	if _, ok := byQN["App.Calc.Inner.Ping"]; !ok {
		t.Errorf("missing App.Calc.Inner.Ping in %+v", funcs)
	}
}
