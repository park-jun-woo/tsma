package index

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsNameField builds a `name` field identifier leaf node.
func rsNameField(text string) *treesitter.Node {
	return &treesitter.Node{Type: "identifier", Field: "name", Text: text}
}

// rsVis builds a `(visibility_modifier)` leaf with the given keyword text.
func rsVis(kw string) *treesitter.Node {
	return &treesitter.Node{Type: "visibility_modifier", Text: kw}
}

// rsFuncNode builds a function_item with an optional visibility_modifier child
// and a name field, spanning rows sRow..eRow (0-based).
func rsFuncNode(name, vis string, sRow, eRow int) *treesitter.Node {
	n := &treesitter.Node{Type: "function_item", SRow: sRow, ERow: eRow}
	if vis != "" {
		n.Children = append(n.Children, rsVis(vis))
	}
	if name != "" {
		n.Children = append(n.Children, rsNameField(name))
	}
	return n
}

// rsCfgTestAttr builds a #[cfg(test)] attribute_item subtree (cfg + test
// identifiers nested under the attribute).
func rsCfgTestAttr() *treesitter.Node {
	return &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
		{Type: "attribute", Children: []*treesitter.Node{
			{Type: "identifier", Text: "cfg"},
			{Type: "token_tree", Children: []*treesitter.Node{
				{Type: "identifier", Text: "all"},
				{Type: "identifier", Text: "test"},
			}},
		}},
	}}
}

// TestRsAttrCfgTest covers cfg+test (true), cfg-only, test-only, and neither
// (false), plus the non-identifier walk branch and the "other identifier" guard.
func TestRsAttrCfgTest(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want bool
	}{
		{"cfg(test) with extra identifiers", rsCfgTestAttr(), true},
		{
			desc: "cfg only -> false",
			node: &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
				{Type: "identifier", Text: "cfg"},
				{Type: "identifier", Text: "feature"},
			}},
			want: false,
		},
		{
			desc: "test only (#[test]) -> false",
			node: &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
				{Type: "identifier", Text: "test"},
			}},
			want: false,
		},
		{
			desc: "neither -> false",
			node: &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
				{Type: "identifier", Text: "derive"},
			}},
			want: false,
		},
	}
	for _, c := range cases {
		if got := rsAttrCfgTest(c.node); got != c.want {
			t.Errorf("%s: rsAttrCfgTest = %v, want %v", c.desc, got, c.want)
		}
	}
}

// TestRsImplTypeName covers nil, plain type_identifier, generic_type with and
// without an inner `type` field, and the default (unexpected shape) branch.
func TestRsImplTypeName(t *testing.T) {
	if got := rsImplTypeName(nil); got != "" {
		t.Errorf("nil: rsImplTypeName = %q, want empty", got)
	}

	plain := &treesitter.Node{Type: "type_identifier", Text: "Foo"}
	if got := rsImplTypeName(plain); got != "Foo" {
		t.Errorf("type_identifier: rsImplTypeName = %q, want Foo", got)
	}

	generic := &treesitter.Node{Type: "generic_type", Children: []*treesitter.Node{
		{Type: "type_identifier", Field: "type", Text: "Bar"},
		{Type: "type_arguments", Text: "<T>"},
	}}
	if got := rsImplTypeName(generic); got != "Bar" {
		t.Errorf("generic_type: rsImplTypeName = %q, want Bar", got)
	}

	genericNoInner := &treesitter.Node{Type: "generic_type"}
	if got := rsImplTypeName(genericNoInner); got != "" {
		t.Errorf("generic_type without inner type: rsImplTypeName = %q, want empty", got)
	}

	other := &treesitter.Node{Type: "scoped_type_identifier", Text: "a::B"}
	if got := rsImplTypeName(other); got != "a::B" {
		t.Errorf("default: rsImplTypeName = %q, want a::B", got)
	}
}

// TestRsNodeExported covers pub, pub(crate), no visibility, a non-pub
// visibility_modifier (Type matches, prefix fails), and a non-visibility child
// carrying "pub" text (Type mismatch).
func TestRsNodeExported(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want bool
	}{
		{"pub", rsFuncNode("a", "pub", 0, 0), true},
		{"pub(crate)", rsFuncNode("a", "pub(crate)", 0, 0), true},
		{"no visibility", rsFuncNode("a", "", 0, 0), false},
		{
			desc: "visibility_modifier without pub prefix",
			node: &treesitter.Node{Type: "function_item", Children: []*treesitter.Node{
				{Type: "visibility_modifier", Text: "crate"}, rsNameField("a"),
			}},
			want: false,
		},
		{
			desc: "pub text on a non-visibility node is ignored",
			node: &treesitter.Node{Type: "function_item", Children: []*treesitter.Node{
				{Type: "identifier", Text: "pub"}, rsNameField("a"),
			}},
			want: false,
		},
	}
	for _, c := range cases {
		if got := rsNodeExported(c.node); got != c.want {
			t.Errorf("%s: rsNodeExported = %v, want %v", c.desc, got, c.want)
		}
	}
}

// TestRsFuncFromNode covers nameless / empty-name (ok=false) and the successful
// pub and non-pub conversions (span, qualified name, exported).
func TestRsFuncFromNode(t *testing.T) {
	scopes := []rsScope{{receiver: "Calc"}}

	if _, ok := rsFuncFromNode(&treesitter.Node{Type: "function_item"}, "src", scopes, "src/lib.rs"); ok {
		t.Error("nameless declaration should return ok=false")
	}

	emptyName := &treesitter.Node{Type: "function_item", Children: []*treesitter.Node{rsNameField("")}}
	if _, ok := rsFuncFromNode(emptyName, "src", scopes, "src/lib.rs"); ok {
		t.Error("empty-name declaration should return ok=false")
	}

	node := rsFuncNode("compute", "pub", 4, 6)
	fn, ok := rsFuncFromNode(node, "src", scopes, "src/lib.rs")
	if !ok {
		t.Fatal("valid declaration should return ok=true")
	}
	if fn.QualifiedName != "src.Calc.compute" || fn.Name != "compute" ||
		fn.File != "src/lib.rs" || fn.StartLine != 5 || fn.EndLine != 7 ||
		!fn.Exported || fn.Status != model.StatusTodo {
		t.Errorf("rsFuncFromNode = %+v", fn)
	}

	priv := rsFuncNode("helper", "", 1, 1)
	pf, ok := rsFuncFromNode(priv, "src", nil, "src/lib.rs")
	if !ok || pf.Exported || pf.QualifiedName != "src.helper" {
		t.Errorf("non-pub: rsFuncFromNode = %+v, ok=%v", pf, ok)
	}
}

// TestEmitRsFunc covers the cfgTestActive (pending and scope) skips, the
// nameless skip, and the successful append.
func TestEmitRsFunc(t *testing.T) {
	// pending guard: skipped.
	var out []model.Function
	emitRsFunc(rsFuncNode("a", "pub", 1, 1), "src", nil, "src/lib.rs", &out, true)
	if len(out) != 0 {
		t.Errorf("pending guard: out = %+v, want empty", out)
	}

	// scope cfgTest guard: skipped.
	out = nil
	emitRsFunc(rsFuncNode("a", "pub", 1, 1), "src", []rsScope{{module: "tests", cfgTest: true}}, "src/lib.rs", &out, false)
	if len(out) != 0 {
		t.Errorf("scope guard: out = %+v, want empty", out)
	}

	// not guarded but nameless: not appended.
	out = nil
	emitRsFunc(&treesitter.Node{Type: "function_item"}, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 0 {
		t.Errorf("nameless: out = %+v, want empty", out)
	}

	// not guarded, valid: appended.
	out = nil
	emitRsFunc(rsFuncNode("go", "pub", 2, 2), "src", nil, "src/lib.rs", &out, false)
	if len(out) != 1 || out[0].Name != "go" {
		t.Errorf("valid: out = %+v, want single go", out)
	}
}

// TestDispatchRsMember covers attribute_item (cfg(test) and non-test, with and
// without prior pending), function_item (emitted and guarded), impl_item,
// mod_item, and the default carry-through branch.
func TestDispatchRsMember(t *testing.T) {
	implNode := func() *treesitter.Node {
		return &treesitter.Node{Type: "impl_item", Children: []*treesitter.Node{
			{Type: "type_identifier", Field: "type", Text: "Foo"},
			{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
				rsFuncNode("m", "pub", 1, 1),
			}},
		}}
	}
	modNode := func() *treesitter.Node {
		return &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{
			rsNameField("inner"),
			{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
				rsFuncNode("g", "pub", 2, 2),
			}},
		}}
	}

	cases := []struct {
		desc        string
		node        *treesitter.Node
		pending     bool
		wantPending bool
		wantFuncs   int
	}{
		{"cfg(test) attribute sets pending", rsCfgTestAttr(), false, true, 0},
		{"non-test attribute keeps pending false", &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "derive"}}}, false, false, 0},
		{"non-test attribute when already pending stays true", &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "derive"}}}, true, true, 0},
		{"function_item emitted", rsFuncNode("f", "pub", 0, 0), false, false, 1},
		{"function_item guarded by pending", rsFuncNode("f", "pub", 0, 0), true, false, 0},
		{"impl_item recurses", implNode(), false, false, 1},
		{"mod_item recurses", modNode(), false, false, 1},
		{"default (use) carries pending", &treesitter.Node{Type: "use_declaration"}, true, true, 0},
		{"default (comment) carries pending false", &treesitter.Node{Type: "line_comment"}, false, false, 0},
	}
	for _, c := range cases {
		var out []model.Function
		gotPending := dispatchRsMember(c.node, "src", nil, "src/lib.rs", &out, c.pending)
		if gotPending != c.wantPending {
			t.Errorf("%s: pending = %v, want %v", c.desc, gotPending, c.wantPending)
		}
		if len(out) != c.wantFuncs {
			t.Errorf("%s: emitted %d funcs (%+v), want %d", c.desc, len(out), out, c.wantFuncs)
		}
	}
}

// TestCollectRsMembers covers sibling guard threading (a #[cfg(test)] attr
// excludes the following fn), a non-test attribute (fn kept), and the empty
// container (loop body never runs).
func TestCollectRsMembers(t *testing.T) {
	// cfg(test) attribute followed by a fn -> fn excluded.
	container := &treesitter.Node{Type: "source_file", Children: []*treesitter.Node{
		rsCfgTestAttr(),
		rsFuncNode("test_only", "pub", 1, 1),
		rsFuncNode("real", "pub", 2, 2),
	}}
	var out []model.Function
	collectRsMembers(container, "src", nil, "src/lib.rs", &out)
	if len(out) != 1 || out[0].Name != "real" {
		t.Fatalf("guard threading: out = %+v, want single real", out)
	}

	// non-test attribute followed by a fn -> fn kept.
	out = nil
	container2 := &treesitter.Node{Type: "source_file", Children: []*treesitter.Node{
		{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "inline"}}},
		rsFuncNode("kept", "pub", 3, 3),
	}}
	collectRsMembers(container2, "src", nil, "src/lib.rs", &out)
	if len(out) != 1 || out[0].Name != "kept" {
		t.Errorf("non-test attribute: out = %+v, want single kept", out)
	}

	// empty container.
	out = nil
	collectRsMembers(&treesitter.Node{Type: "source_file"}, "src", nil, "src/lib.rs", &out)
	if len(out) != 0 {
		t.Errorf("empty container: out = %+v, want empty", out)
	}
}

// TestWalkRsImpl covers the named-type recurse, the nameless impl (empty
// receiver, still descends), the body-less impl (returns), and the cfgTest
// inheritance (methods of a #[cfg(test)] impl excluded).
func TestWalkRsImpl(t *testing.T) {
	body := func() *treesitter.Node {
		return &treesitter.Node{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			rsFuncNode("method", "pub", 1, 1),
		}}
	}

	// named impl with a body method.
	var out []model.Function
	named := &treesitter.Node{Type: "impl_item", Children: []*treesitter.Node{
		{Type: "type_identifier", Field: "type", Text: "Foo"}, body(),
	}}
	walkRsImpl(named, "src", []rsScope{{module: "util"}}, "src/lib.rs", &out, false)
	if len(out) != 1 || out[0].QualifiedName != "src.util::Foo.method" {
		t.Errorf("named impl: out = %+v, want src.util::Foo.method", out)
	}

	// nameless impl: empty receiver, still descends.
	out = nil
	nameless := &treesitter.Node{Type: "impl_item", Children: []*treesitter.Node{body()}}
	walkRsImpl(nameless, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 1 || out[0].QualifiedName != "src.method" {
		t.Errorf("nameless impl: out = %+v, want src.method", out)
	}

	// body-less impl: returns, nothing emitted.
	out = nil
	noBody := &treesitter.Node{Type: "impl_item", Children: []*treesitter.Node{
		{Type: "type_identifier", Field: "type", Text: "Foo"},
	}}
	walkRsImpl(noBody, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 0 {
		t.Errorf("body-less impl: out = %+v, want empty", out)
	}

	// cfgTest inheritance: methods excluded.
	out = nil
	walkRsImpl(named, "src", nil, "src/lib.rs", &out, true)
	if len(out) != 0 {
		t.Errorf("cfgTest impl: out = %+v, want empty", out)
	}
}

// TestWalkRsMod covers the named-mod recurse, the nameless mod (skipped), the
// body-less `mod foo;` declaration (skipped), and the cfgTest inheritance.
func TestWalkRsMod(t *testing.T) {
	body := func() *treesitter.Node {
		return &treesitter.Node{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			rsFuncNode("inner_fn", "pub", 1, 1),
		}}
	}

	// named mod with a body fn.
	var out []model.Function
	named := &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{
		rsNameField("nested"), body(),
	}}
	walkRsMod(named, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 1 || out[0].QualifiedName != "src.nested.inner_fn" {
		t.Errorf("named mod: out = %+v, want src.nested.inner_fn", out)
	}

	// nameless mod: skipped.
	out = nil
	walkRsMod(&treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{body()}}, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 0 {
		t.Errorf("nameless mod: out = %+v, want empty", out)
	}

	// body-less `mod foo;`: skipped.
	out = nil
	walkRsMod(&treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{rsNameField("foo")}}, "src", nil, "src/lib.rs", &out, false)
	if len(out) != 0 {
		t.Errorf("body-less mod: out = %+v, want empty", out)
	}

	// cfgTest inheritance: fn excluded.
	out = nil
	walkRsMod(named, "src", nil, "src/lib.rs", &out, true)
	if len(out) != 0 {
		t.Errorf("cfgTest mod: out = %+v, want empty", out)
	}
}

// TestExtractRustFunctions covers the nil-root guard and a full extraction over
// a hand-built source_file: a free fn, an impl method, a nested-mod fn, and a
// #[cfg(test)] mod whose functions must NOT be indexed.
func TestExtractRustFunctions(t *testing.T) {
	if got := extractRustFunctions(nil, "src/lib.rs", "src"); got != nil {
		t.Errorf("nil root: extractRustFunctions = %+v, want nil", got)
	}

	implMethod := &treesitter.Node{Type: "impl_item", Children: []*treesitter.Node{
		{Type: "type_identifier", Field: "type", Text: "Calc"},
		{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			rsFuncNode("new", "pub", 10, 11),
		}},
	}}
	nestedMod := &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{
		rsNameField("nested"),
		{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			rsFuncNode("double", "pub", 20, 21),
		}},
	}}
	testMod := &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{
		rsNameField("tests"),
		{Type: "declaration_list", Field: "body", Children: []*treesitter.Node{
			rsFuncNode("it_works", "", 30, 31),
		}},
	}}
	root := &treesitter.Node{Type: "source_file", Children: []*treesitter.Node{
		rsFuncNode("add", "pub", 0, 1),
		implMethod,
		nestedMod,
		rsCfgTestAttr(),
		testMod, // guarded by the preceding cfg(test) attribute
	}}

	funcs := extractRustFunctions(root, "src/lib.rs", "src")
	byQN := map[string]model.Function{}
	for _, f := range funcs {
		byQN[f.QualifiedName] = f
		if f.Name == "it_works" {
			t.Errorf("#[cfg(test)] function should not be indexed: %+v", f)
		}
	}
	for _, qn := range []string{"src.add", "src.Calc.new", "src.nested.double"} {
		if _, ok := byQN[qn]; !ok {
			t.Errorf("missing %q in %+v", qn, funcs)
		}
	}
}
