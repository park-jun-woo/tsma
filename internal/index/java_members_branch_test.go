package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// nameField builds a `name` field identifier leaf node.
func nameField(text string) *treesitter.Node {
	return &treesitter.Node{Type: "identifier", Field: "name", Text: text}
}

// methodNode builds a method_declaration with an optional modifiers child and a
// name field, spanning rows sRow..eRow (0-based).
func methodNode(name string, mods string, sRow, eRow int) *treesitter.Node {
	n := &treesitter.Node{Type: "method_declaration", SRow: sRow, ERow: eRow}
	if mods != "" {
		n.Children = append(n.Children, &treesitter.Node{Type: "modifiers", Text: mods})
	}
	if name != "" {
		n.Children = append(n.Children, nameField(name))
	}
	return n
}

// TestJavaNodeExported covers all three modifier shapes.
func TestJavaNodeExported(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want bool
	}{
		{
			desc: "public modifier among several",
			node: &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{
				{Type: "modifiers", Text: "public static final"},
			}},
			want: true,
		},
		{
			desc: "modifiers without public",
			node: &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{
				{Type: "modifiers", Text: "private static"},
			}},
			want: false,
		},
		{
			desc: "no modifiers child",
			node: &treesitter.Node{Type: "method_declaration"},
			want: false,
		},
	}
	for _, c := range cases {
		if got := javaNodeExported(c.node); got != c.want {
			t.Errorf("%s: javaNodeExported = %v, want %v", c.desc, got, c.want)
		}
	}
}

// TestJavaFuncFromNode covers the nameless / empty-name false branches and the
// successful conversion (span, qualified name, exported).
func TestJavaFuncFromNode(t *testing.T) {
	scopes := []javaScope{{typeName: "Outer"}}

	// nameless: no `name` field at all.
	if _, ok := javaFuncFromNode(&treesitter.Node{Type: "method_declaration"}, "pkg", scopes, "F.java"); ok {
		t.Error("nameless declaration should return ok=false")
	}

	// empty name: `name` field present but blank text.
	emptyName := &treesitter.Node{Type: "method_declaration", Children: []*treesitter.Node{nameField("")}}
	if _, ok := javaFuncFromNode(emptyName, "pkg", scopes, "F.java"); ok {
		t.Error("empty-name declaration should return ok=false")
	}

	// valid public method.
	node := methodNode("compute", "public", 4, 6)
	fn, ok := javaFuncFromNode(node, "pkg", scopes, "F.java")
	if !ok {
		t.Fatal("valid declaration should return ok=true")
	}
	if fn.QualifiedName != "pkg.Outer.compute" || fn.Name != "compute" ||
		fn.File != "F.java" || fn.StartLine != 5 || fn.EndLine != 7 ||
		!fn.Exported || fn.Status != model.StatusTodo {
		t.Errorf("javaFuncFromNode = %+v", fn)
	}
}

// TestWalkJavaTypeDecl covers the nameless, empty-name, body-less and full
// (recurse) branches.
func TestWalkJavaTypeDecl(t *testing.T) {
	// nameless type: skipped.
	var out []model.Function
	walkJavaTypeDecl(&treesitter.Node{Type: "class_declaration"}, "pkg", nil, "F.java", &out)
	if len(out) != 0 {
		t.Errorf("nameless type: out = %+v, want empty", out)
	}

	// empty-name type: skipped.
	out = nil
	emptyName := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{nameField("")}}
	walkJavaTypeDecl(emptyName, "pkg", nil, "F.java", &out)
	if len(out) != 0 {
		t.Errorf("empty-name type: out = %+v, want empty", out)
	}

	// named type with no body: returns after pushing scope, no funcs.
	out = nil
	noBody := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{nameField("NoBody")}}
	walkJavaTypeDecl(noBody, "pkg", nil, "F.java", &out)
	if len(out) != 0 {
		t.Errorf("body-less type: out = %+v, want empty", out)
	}

	// named type with a body containing one method: recurse and qualify.
	out = nil
	body := &treesitter.Node{Type: "class_body", Field: "body", Children: []*treesitter.Node{
		methodNode("m", "public", 1, 1),
	}}
	withBody := &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{
		nameField("Widget"), body,
	}}
	walkJavaTypeDecl(withBody, "pkg", []javaScope{{typeName: "Top"}}, "F.java", &out)
	if len(out) != 1 || out[0].QualifiedName != "pkg.Top.Widget.m" {
		t.Errorf("named-with-body: out = %+v, want single pkg.Top.Widget.m", out)
	}
}

// TestDispatchJavaMember covers every recognized member shape plus the ignored
// default case.
func TestDispatchJavaMember(t *testing.T) {
	methodChild := func() *treesitter.Node { return methodNode("m", "public", 1, 1) }

	cases := []struct {
		desc string
		node *treesitter.Node
		want int // number of funcs emitted
	}{
		{"method_declaration", methodNode("am", "", 1, 1), 1},
		{"constructor_declaration", &treesitter.Node{Type: "constructor_declaration", Children: []*treesitter.Node{nameField("Ctor")}}, 1},
		{"compact_constructor_declaration", &treesitter.Node{Type: "compact_constructor_declaration", Children: []*treesitter.Node{nameField("Rec")}}, 1},
		{"class_declaration recurses", &treesitter.Node{Type: "class_declaration", Children: []*treesitter.Node{
			nameField("C"),
			{Type: "class_body", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}, 1},
		{"interface_declaration recurses", &treesitter.Node{Type: "interface_declaration", Children: []*treesitter.Node{
			nameField("I"),
			{Type: "interface_body", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}, 1},
		{"enum_declaration recurses", &treesitter.Node{Type: "enum_declaration", Children: []*treesitter.Node{
			nameField("E"),
			{Type: "enum_body", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}, 1},
		{"record_declaration recurses", &treesitter.Node{Type: "record_declaration", Children: []*treesitter.Node{
			nameField("R"),
			{Type: "class_body", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}, 1},
		{"annotation_type_declaration recurses", &treesitter.Node{Type: "annotation_type_declaration", Children: []*treesitter.Node{
			nameField("A"),
			{Type: "annotation_type_body", Field: "body", Children: []*treesitter.Node{methodChild()}},
		}}, 1},
		{"enum_body_declarations descends without scope", &treesitter.Node{Type: "enum_body_declarations", Children: []*treesitter.Node{methodChild()}}, 1},
		{"field_declaration ignored", &treesitter.Node{Type: "field_declaration"}, 0},
	}
	for _, c := range cases {
		var out []model.Function
		dispatchJavaMember(c.node, "pkg", nil, "F.java", &out)
		if len(out) != c.want {
			t.Errorf("%s: emitted %d funcs (%+v), want %d", c.desc, len(out), out, c.want)
		}
	}
}

// TestCollectJavaMembers covers the iteration over a container's children.
func TestCollectJavaMembers(t *testing.T) {
	container := &treesitter.Node{Type: "class_body", Children: []*treesitter.Node{
		methodNode("a", "", 1, 1),
		methodNode("b", "", 2, 2),
		{Type: "field_declaration"}, // ignored
	}}
	var out []model.Function
	collectJavaMembers(container, "pkg", nil, "F.java", &out)
	if len(out) != 2 {
		t.Fatalf("collectJavaMembers emitted %d funcs (%+v), want 2", len(out), out)
	}

	// empty container: loop body never runs.
	out = nil
	collectJavaMembers(&treesitter.Node{Type: "class_body"}, "pkg", nil, "F.java", &out)
	if len(out) != 0 {
		t.Errorf("empty container: out = %+v, want empty", out)
	}
}

// TestJavaPackageNameBranches covers nil-declaration, dotted package, an
// empty-text identifier (Text != "" false), and a non-identifier node (Type
// false) inside the package declaration.
func TestJavaPackageNameBranches(t *testing.T) {
	// no package_declaration child -> "".
	root := &treesitter.Node{Type: "program", Children: []*treesitter.Node{
		{Type: "class_declaration"},
	}}
	if got := javaPackageName(root); got != "" {
		t.Errorf("default package: javaPackageName = %q, want empty", got)
	}

	// dotted package with a nested scoped_identifier (non-identifier node) and
	// one stray empty-text identifier that must be skipped.
	pkgDecl := &treesitter.Node{Type: "package_declaration", Children: []*treesitter.Node{
		{Type: "scoped_identifier", Children: []*treesitter.Node{
			{Type: "identifier", Text: "com"},
			{Type: "identifier", Text: "example"},
			{Type: "identifier", Text: ""}, // empty -> skipped
			{Type: "identifier", Text: "app"},
		}},
	}}
	root2 := &treesitter.Node{Type: "program", Children: []*treesitter.Node{pkgDecl}}
	if got := javaPackageName(root2); got != "com.example.app" {
		t.Errorf("dotted package: javaPackageName = %q, want com.example.app", got)
	}
}
