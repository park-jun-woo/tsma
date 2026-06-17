package smell

import (
	"os"
	"os/exec"
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// tn builds a treesitter.Node with a field name and children.
func tn(typ, field, text string, children ...*treesitter.Node) *treesitter.Node {
	return &treesitter.Node{Type: typ, Field: field, Text: text, SRow: 4, Children: children}
}

// asAnyMember builds `(x as any).<prop>` as a member_expression node. If prop is
// "" the property field is omitted (exercising the prop=="" branch).
func asAnyMember(prop string) *treesitter.Node {
	inner := tn("as_expression", "", "",
		tn("identifier", "expression", "x"),
		tn("predefined_type", "type", "any"),
	)
	obj := tn("parenthesized_expression", "object", "", inner)
	children := []*treesitter.Node{obj}
	if prop != "" {
		children = append(children, tn("property_identifier", "property", prop))
	}
	return tn("member_expression", "", "", children...)
}

func TestScanTSUnavailable(t *testing.T) {
	// a bogus tree-sitter command LookPath cannot resolve -> ScanTS errors out
	// (the command=="" guard), so the caller can ignore the file.
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	if got, err := ScanTS("whatever.test.ts"); err == nil || got != nil {
		t.Errorf("ScanTS(unavailable) = (%v,%v), want (nil,err)", got, err)
	}
}

func TestScanTSParseError(t *testing.T) {
	cmd := os.Getenv("TSMA_TREE_SITTER")
	if cmd == "" {
		if _, err := exec.LookPath("tree-sitter"); err != nil {
			t.Skip("tree-sitter unavailable")
		}
	}
	// a non-existent file: the CLI is present so command!="" but ParseFile fails.
	if got, err := ScanTS("../../testdata/typescript/no-such-file.test.ts"); err == nil || got != nil {
		t.Errorf("ScanTS(missing file) = (%v,%v), want (nil,err)", got, err)
	}
}

func TestUnwrapAsExpr(t *testing.T) {
	if unwrapAsExpr(nil) != nil {
		t.Error("unwrapAsExpr(nil) != nil")
	}
	// parenthesized_expression -> inner as_expression.
	ae := tn("as_expression", "", "")
	paren := tn("parenthesized_expression", "object", "", ae)
	if got := unwrapAsExpr(paren); got != ae {
		t.Errorf("unwrapAsExpr(paren) = %+v, want inner as_expression", got)
	}
	// parenthesized_expression with no as_expression -> nil.
	parenEmpty := tn("parenthesized_expression", "object", "", tn("identifier", "", "x"))
	if got := unwrapAsExpr(parenEmpty); got != nil {
		t.Errorf("unwrapAsExpr(paren w/o cast) = %+v, want nil", got)
	}
	// direct as_expression -> itself.
	direct := tn("as_expression", "object", "")
	if got := unwrapAsExpr(direct); got != direct {
		t.Errorf("unwrapAsExpr(as_expression) = %+v, want itself", got)
	}
	// anything else -> nil.
	if got := unwrapAsExpr(tn("identifier", "object", "x")); got != nil {
		t.Errorf("unwrapAsExpr(identifier) = %+v, want nil", got)
	}
}

func TestAsExprIsAny(t *testing.T) {
	// cast to any -> true.
	any := tn("as_expression", "", "", tn("predefined_type", "type", "any"))
	if !asExprIsAny(any) {
		t.Error("asExprIsAny(any) = false")
	}
	// cast to a different predefined type -> false.
	str := tn("as_expression", "", "", tn("predefined_type", "type", "string"))
	if asExprIsAny(str) {
		t.Error("asExprIsAny(string) = true")
	}
	// cast to a named (non-predefined) type -> false.
	named := tn("as_expression", "", "", tn("type_identifier", "type", "Foo"))
	if asExprIsAny(named) {
		t.Error("asExprIsAny(named) = true")
	}
}

func TestDetectTSAsAny(t *testing.T) {
	root := tn("program", "", "",
		// fires: (x as any).secret
		asAnyMember("secret"),
		// fires with empty prop note: (x as any) member with no property field.
		asAnyMember(""),
		// not a member_expression -> ignored.
		tn("call_expression", "", ""),
		// member_expression whose object is not an as-cast -> ignored.
		tn("member_expression", "", "",
			tn("identifier", "object", "obj"),
			tn("property_identifier", "property", "field"),
		),
		// member_expression cast to a named type -> ignored.
		tn("member_expression", "", "",
			tn("parenthesized_expression", "object", "",
				tn("as_expression", "", "", tn("type_identifier", "type", "Foo")),
			),
			tn("property_identifier", "property", "field"),
		),
	)
	got := detectTSAsAny(root, "f.ts")
	if len(got) != 2 {
		t.Fatalf("detectTSAsAny fired %d times, want 2: %+v", len(got), got)
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-TS-001" || f.File != "f.ts" || f.Line != 5 {
			t.Errorf("finding = %+v", f)
		}
	}
	if got[0].Note != "(as any).secret" {
		t.Errorf("note = %q", got[0].Note)
	}
	if got[1].Note != "(as any)." {
		t.Errorf("empty-prop note = %q", got[1].Note)
	}
}

func TestDetectTSReflect(t *testing.T) {
	root := tn("program", "", "",
		// fires: Reflect.get
		tn("member_expression", "", "",
			tn("identifier", "object", "Reflect"),
			tn("property_identifier", "property", "get"),
		),
		// fires with empty prop: Reflect member with no property field.
		tn("member_expression", "", "",
			tn("identifier", "object", "Reflect"),
		),
		// not a member_expression -> ignored.
		tn("identifier", "", "Reflect"),
		// object missing -> ignored.
		tn("member_expression", "", "",
			tn("property_identifier", "property", "get"),
		),
		// object not an identifier -> ignored.
		tn("member_expression", "", "",
			tn("member_expression", "object", ""),
			tn("property_identifier", "property", "get"),
		),
		// object identifier but not Reflect -> ignored.
		tn("member_expression", "", "",
			tn("identifier", "object", "Math"),
			tn("property_identifier", "property", "max"),
		),
	)
	got := detectTSReflect(root, "f.ts")
	if len(got) != 2 {
		t.Fatalf("detectTSReflect fired %d times, want 2: %+v", len(got), got)
	}
	if got[0].Rule != "TS-REFL-TS-002" || got[0].Note != "Reflect.get" {
		t.Errorf("finding[0] = %+v", got[0])
	}
	if got[1].Note != "Reflect." {
		t.Errorf("empty-prop note = %q", got[1].Note)
	}
}

func TestDetectTSOwnProperty(t *testing.T) {
	root := tn("program", "", "",
		// fires: Object.getOwnPropertyNames
		tn("member_expression", "", "",
			tn("identifier", "object", "Object"),
			tn("property_identifier", "property", "getOwnPropertyNames"),
		),
		// fires: Object.getOwnPropertyDescriptor
		tn("member_expression", "", "",
			tn("identifier", "object", "Object"),
			tn("property_identifier", "property", "getOwnPropertyDescriptor"),
		),
		// not a member_expression -> ignored.
		tn("identifier", "", "Object"),
		// object missing -> ignored.
		tn("member_expression", "", "",
			tn("property_identifier", "property", "getOwnPropertyNames"),
		),
		// object not identifier -> ignored.
		tn("member_expression", "", "",
			tn("call_expression", "object", ""),
			tn("property_identifier", "property", "getOwnPropertyNames"),
		),
		// object not "Object" -> ignored.
		tn("member_expression", "", "",
			tn("identifier", "object", "Reflect"),
			tn("property_identifier", "property", "getOwnPropertyNames"),
		),
		// property missing -> ignored.
		tn("member_expression", "", "",
			tn("identifier", "object", "Object"),
		),
		// public iteration method (not in reach map) -> ignored.
		tn("member_expression", "", "",
			tn("identifier", "object", "Object"),
			tn("property_identifier", "property", "keys"),
		),
	)
	got := detectTSOwnProperty(root, "f.ts")
	if len(got) != 2 {
		t.Fatalf("detectTSOwnProperty fired %d times, want 2: %+v", len(got), got)
	}
	if got[0].Rule != "TS-REFL-TS-003" || got[0].Note != "Object.getOwnPropertyNames" {
		t.Errorf("finding[0] = %+v", got[0])
	}
	if got[1].Note != "Object.getOwnPropertyDescriptor" {
		t.Errorf("finding[1] = %+v", got[1])
	}
}
