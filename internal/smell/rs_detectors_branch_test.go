package smell

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsScopedRead builds a std::ptr::read scoped_identifier call function node.
// scoped_identifier{ path: scoped_identifier{ name=ptr }, name=read }.
func rsScopedRead(method, seg string) *treesitter.Node {
	return &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
		{Type: "scoped_identifier", Field: "path", Children: []*treesitter.Node{
			{Type: "identifier", Field: "path", Text: "std"},
			{Type: "identifier", Field: "name", Text: seg},
		}},
		{Type: "identifier", Field: "name", Text: method},
	}}
}

// TestRsAttrMentionsTest covers the test-mentioning (true) and absent (false)
// attribute_item subtrees, plus the non-identifier walk node.
func TestRsAttrMentionsTest(t *testing.T) {
	withTest := &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
		{Type: "token_tree", Children: []*treesitter.Node{{Type: "identifier", Text: "test"}}},
	}}
	if !rsAttrMentionsTest(withTest) {
		t.Error("#[...test...]: rsAttrMentionsTest = false, want true")
	}
	without := &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
		{Type: "identifier", Text: "derive"},
	}}
	if rsAttrMentionsTest(without) {
		t.Error("#[derive]: rsAttrMentionsTest = true, want false")
	}
}

// TestRsCallName covers nil, identifier, scoped_identifier (with/without name),
// and the default branch.
func TestRsCallName(t *testing.T) {
	cases := []struct {
		desc string
		node *treesitter.Node
		want string
	}{
		{"nil", nil, ""},
		{"identifier", &treesitter.Node{Type: "identifier", Text: "transmute"}, "transmute"},
		{
			desc: "scoped_identifier with name",
			node: &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
				{Type: "identifier", Field: "name", Text: "transmute"},
			}},
			want: "transmute",
		},
		{"scoped_identifier without name", &treesitter.Node{Type: "scoped_identifier"}, ""},
		{"default", &treesitter.Node{Type: "field_expression"}, ""},
	}
	for _, c := range cases {
		if got := rsCallName(c.node); got != c.want {
			t.Errorf("%s: rsCallName = %q, want %q", c.desc, got, c.want)
		}
	}
}

// TestRsFnIsUnsafe covers the unsafe modifier (true), a non-unsafe modifier set
// (false), and the no-modifiers branch.
func TestRsFnIsUnsafe(t *testing.T) {
	unsafeFn := &treesitter.Node{Type: "function_item", Children: []*treesitter.Node{
		{Type: "function_modifiers", Text: "unsafe"},
	}}
	if !rsFnIsUnsafe(unsafeFn) {
		t.Error("unsafe fn: rsFnIsUnsafe = false, want true")
	}
	asyncFn := &treesitter.Node{Type: "function_item", Children: []*treesitter.Node{
		{Type: "function_modifiers", Text: "async"},
	}}
	if rsFnIsUnsafe(asyncFn) {
		t.Error("async fn: rsFnIsUnsafe = true, want false")
	}
	plainFn := &treesitter.Node{Type: "function_item"}
	if rsFnIsUnsafe(plainFn) {
		t.Error("plain fn: rsFnIsUnsafe = true, want false")
	}
}

// TestRsPtrFieldAccess covers nil, as_ptr, as_mut_ptr, and a non-pointer field.
func TestRsPtrFieldAccess(t *testing.T) {
	if _, ok := rsPtrFieldAccess(nil); ok {
		t.Error("nil: rsPtrFieldAccess ok = true, want false")
	}
	for _, name := range []string{"as_ptr", "as_mut_ptr"} {
		note, ok := rsPtrFieldAccess(&treesitter.Node{Type: "field_identifier", Text: name})
		if !ok || note != name+"()" {
			t.Errorf("%s: rsPtrFieldAccess = (%q,%v), want (%q,true)", name, note, ok, name+"()")
		}
	}
	if _, ok := rsPtrFieldAccess(&treesitter.Node{Type: "field_identifier", Text: "len"}); ok {
		t.Error("len: rsPtrFieldAccess ok = true, want false")
	}
}

// TestRsPtrScopedCall covers nil, wrong type, missing/unknown name, missing path,
// non-scoped path, wrong/missing path segment, and the full std::ptr::read match.
func TestRsPtrScopedCall(t *testing.T) {
	// full match.
	note, ok := rsPtrScopedCall(rsScopedRead("read", "ptr"))
	if !ok || note != "ptr::read()" {
		t.Errorf("std::ptr::read: rsPtrScopedCall = (%q,%v), want (ptr::read(),true)", note, ok)
	}

	bad := []struct {
		desc string
		node *treesitter.Node
	}{
		{"nil", nil},
		{"wrong type", &treesitter.Node{Type: "identifier", Text: "read"}},
		{"missing name", &treesitter.Node{Type: "scoped_identifier"}},
		{"unknown name", rsScopedRead("clone", "ptr")},
		{
			desc: "missing path",
			node: &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
				{Type: "identifier", Field: "name", Text: "read"},
			}},
		},
		{
			desc: "non-scoped path",
			node: &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
				{Type: "identifier", Field: "path", Text: "reader"},
				{Type: "identifier", Field: "name", Text: "read"},
			}},
		},
		{"wrong path segment", rsScopedRead("read", "mem")},
		{
			desc: "path without name segment",
			node: &treesitter.Node{Type: "scoped_identifier", Children: []*treesitter.Node{
				{Type: "scoped_identifier", Field: "path"},
				{Type: "identifier", Field: "name", Text: "read"},
			}},
		},
	}
	for _, c := range bad {
		if _, ok := rsPtrScopedCall(c.node); ok {
			t.Errorf("%s: rsPtrScopedCall ok = true, want false", c.desc)
		}
	}
}

// TestRsCollectTestScope covers the attribute (test / non-test), mod_item
// (pending+body, pending+no-body, no-pending), function_item (pending /
// no-pending), and default branches.
func TestRsCollectTestScope(t *testing.T) {
	testAttr := &treesitter.Node{Type: "attribute_item", Children: []*treesitter.Node{
		{Type: "identifier", Text: "test"},
	}}
	body := &treesitter.Node{Type: "declaration_list", Field: "body"}
	modWithBody := &treesitter.Node{Type: "mod_item", Children: []*treesitter.Node{body}}
	fn := &treesitter.Node{Type: "function_item"}

	if p := rsCollectTestScope(testAttr, false, &[]*treesitter.Node{}); !p {
		t.Error("#[test] attr: pending = false, want true")
	}
	if p := rsCollectTestScope(&treesitter.Node{Type: "attribute_item"}, false, &[]*treesitter.Node{}); p {
		t.Error("non-test attr: pending = true, want false")
	}

	// guarded mod -> body appended.
	var out []*treesitter.Node
	if p := rsCollectTestScope(modWithBody, true, &out); p {
		t.Error("guarded mod: pending = true, want false")
	}
	if len(out) != 1 || out[0] != body {
		t.Errorf("guarded mod: out = %v, want [body]", out)
	}
	// guarded mod, no body.
	out = nil
	rsCollectTestScope(&treesitter.Node{Type: "mod_item"}, true, &out)
	if len(out) != 0 {
		t.Errorf("body-less mod: out = %v, want empty", out)
	}
	// unguarded mod.
	out = nil
	rsCollectTestScope(modWithBody, false, &out)
	if len(out) != 0 {
		t.Errorf("unguarded mod: out = %v, want empty", out)
	}
	// guarded function -> appended.
	out = nil
	if p := rsCollectTestScope(fn, true, &out); p {
		t.Error("guarded fn: pending = true, want false")
	}
	if len(out) != 1 || out[0] != fn {
		t.Errorf("guarded fn: out = %v, want [fn]", out)
	}
	// unguarded function.
	out = nil
	rsCollectTestScope(fn, false, &out)
	if len(out) != 0 {
		t.Errorf("unguarded fn: out = %v, want empty", out)
	}
	// default carries pending.
	if p := rsCollectTestScope(&treesitter.Node{Type: "use_declaration"}, true, &out); !p {
		t.Error("default: pending = false, want true (carried)")
	}
}

// TestRsTestScopeNodes covers nil root, a #[cfg(test)] mod body, a top-level
// #[test] function, and an empty root.
func TestRsTestScopeNodes(t *testing.T) {
	if got := rsTestScopeNodes(nil); got != nil {
		t.Errorf("nil root: rsTestScopeNodes = %v, want nil", got)
	}

	body := &treesitter.Node{Type: "declaration_list", Field: "body"}
	fn := &treesitter.Node{Type: "function_item"}
	root := &treesitter.Node{Type: "source_file", Children: []*treesitter.Node{
		{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "test"}}},
		{Type: "mod_item", Children: []*treesitter.Node{body}},
		{Type: "attribute_item", Children: []*treesitter.Node{{Type: "identifier", Text: "test"}}},
		fn,
	}}
	got := rsTestScopeNodes(root)
	if len(got) != 2 || got[0] != body || got[1] != fn {
		t.Errorf("rsTestScopeNodes = %v, want [body fn]", got)
	}

	if got := rsTestScopeNodes(&treesitter.Node{Type: "source_file"}); got != nil {
		t.Errorf("empty root: rsTestScopeNodes = %v, want nil", got)
	}
}

// TestDetectRsUnsafe covers the unsafe_block, unsafe fn, and clean (no finding)
// cases.
func TestDetectRsUnsafe(t *testing.T) {
	scope := &treesitter.Node{Type: "block", Children: []*treesitter.Node{
		{Type: "unsafe_block", SRow: 4},
		{Type: "function_item", SRow: 9, Children: []*treesitter.Node{
			{Type: "function_modifiers", Text: "unsafe"},
		}},
		{Type: "function_item", SRow: 12}, // safe fn -> no finding
		{Type: "let_declaration"},
	}}
	findings := detectRsUnsafe(scope, "f.rs")
	if len(findings) != 2 {
		t.Fatalf("detectRsUnsafe = %+v, want 2 findings", findings)
	}
	if findings[0].Rule != "TS-REFL-RS-001" || findings[0].Note != "unsafe block" || findings[0].Line != 5 {
		t.Errorf("findings[0] = %+v, want unsafe block @5", findings[0])
	}
	if findings[1].Note != "unsafe fn" || findings[1].Line != 10 {
		t.Errorf("findings[1] = %+v, want unsafe fn @10", findings[1])
	}

	// clean scope.
	if got := detectRsUnsafe(&treesitter.Node{Type: "block"}, "f.rs"); got != nil {
		t.Errorf("clean: detectRsUnsafe = %+v, want nil", got)
	}
}

// TestDetectRsTransmute covers a transmute call (finding), a non-transmute call
// (none), and a non-call node (the early return-true branch).
func TestDetectRsTransmute(t *testing.T) {
	scope := &treesitter.Node{Type: "block", Children: []*treesitter.Node{
		{Type: "call_expression", SRow: 6, Children: []*treesitter.Node{
			{Type: "scoped_identifier", Field: "function", Children: []*treesitter.Node{
				{Type: "identifier", Field: "name", Text: "transmute"},
			}},
		}},
		{Type: "call_expression", Children: []*treesitter.Node{
			{Type: "identifier", Field: "function", Text: "to_bytes"},
		}},
		{Type: "let_declaration"},
	}}
	findings := detectRsTransmute(scope, "f.rs")
	if len(findings) != 1 || findings[0].Rule != "TS-REFL-RS-002" || findings[0].Note != "transmute()" || findings[0].Line != 7 {
		t.Fatalf("detectRsTransmute = %+v, want one transmute @7", findings)
	}

	if got := detectRsTransmute(&treesitter.Node{Type: "block"}, "f.rs"); got != nil {
		t.Errorf("clean: detectRsTransmute = %+v, want nil", got)
	}
}

// TestDetectRsPtr covers a std::ptr scoped call, an as_ptr field access, and the
// non-matching call/field/other branches.
func TestDetectRsPtr(t *testing.T) {
	scope := &treesitter.Node{Type: "block", Children: []*treesitter.Node{
		{Type: "call_expression", SRow: 2, Children: []*treesitter.Node{
			func() *treesitter.Node { n := rsScopedRead("read", "ptr"); n.Field = "function"; return n }(),
		}},
		{Type: "field_expression", SRow: 5, Children: []*treesitter.Node{
			{Type: "field_identifier", Field: "field", Text: "as_ptr"},
		}},
		// non-matching call.
		{Type: "call_expression", Children: []*treesitter.Node{
			{Type: "identifier", Field: "function", Text: "read"},
		}},
		// non-matching field.
		{Type: "field_expression", Children: []*treesitter.Node{
			{Type: "field_identifier", Field: "field", Text: "len"},
		}},
		{Type: "let_declaration"},
	}}
	findings := detectRsPtr(scope, "f.rs")
	if len(findings) != 2 {
		t.Fatalf("detectRsPtr = %+v, want 2 findings", findings)
	}
	if findings[0].Note != "ptr::read()" || findings[0].Line != 3 {
		t.Errorf("findings[0] = %+v, want ptr::read() @3", findings[0])
	}
	if findings[1].Note != "as_ptr()" || findings[1].Line != 6 {
		t.Errorf("findings[1] = %+v, want as_ptr() @6", findings[1])
	}

	if got := detectRsPtr(&treesitter.Node{Type: "block"}, "f.rs"); got != nil {
		t.Errorf("clean: detectRsPtr = %+v, want nil", got)
	}
}

// TestScanRsErrors covers the tree-sitter-unavailable and parse-failure error
// branches (the success path is covered by the cheese/calc integration tests).
func TestScanRsErrors(t *testing.T) {
	// tree-sitter unavailable: a bogus absolute override makes ResolveCommand "".
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/tsma-bogus-tree-sitter")
	if _, err := ScanRs("../../testdata/rust/src/cheese.rs"); err == nil {
		t.Error("no tree-sitter: ScanRs err = nil, want error")
	}
}

// TestScanRsParseError covers the ParseFile-error branch: a real CLI over a
// missing file fails to parse.
func TestScanRsParseError(t *testing.T) {
	if treesitter.ResolveCommand() == "" || treesitter.ResolveGrammar("rust") == "" {
		t.Skip("tree-sitter CLI + rust grammar not available")
	}
	if _, err := ScanRs("../../testdata/rust/src/does-not-exist.rs"); err == nil {
		t.Error("missing file: ScanRs err = nil, want parse error")
	}
}
