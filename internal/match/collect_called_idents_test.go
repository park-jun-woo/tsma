package match

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// parseFuncBody parses src (a full file) and returns the body of the first
// top-level function declaration named name.
func parseFuncBody(t *testing.T, src, name string) *ast.BlockStmt {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return fd.Body
		}
	}
	t.Fatalf("func %q not found", name)
	return nil
}

func hasAll(set map[string]struct{}, names ...string) bool {
	for _, n := range names {
		if _, ok := set[n]; !ok {
			return false
		}
	}
	return true
}

func TestCollectCalledIdentsBasic(t *testing.T) {
	body := parseFuncBody(t, `package p
func f() {
	Foo()
	x.Bar()
	if true { Baz() }
}
`, "f")
	got := collectCalledIdents(body)
	if !hasAll(got, "Foo", "Bar", "Baz") {
		t.Errorf("collectCalledIdents = %v, want Foo,Bar,Baz", got)
	}
	if len(got) != 3 {
		t.Errorf("len = %d, want 3 (set dedupe): %v", len(got), got)
	}
}

func TestCollectCalledIdentsDedupes(t *testing.T) {
	body := parseFuncBody(t, `package p
func f() {
	Foo()
	Foo()
	Foo()
}
`, "f")
	got := collectCalledIdents(body)
	if len(got) != 1 || !hasAll(got, "Foo") {
		t.Errorf("collectCalledIdents = %v, want a single Foo", got)
	}
}

func TestCollectCalledIdentsNilBody(t *testing.T) {
	got := collectCalledIdents(nil)
	if got == nil {
		t.Fatal("collectCalledIdents(nil) returned nil map, want empty non-nil")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestCollectCalledIdentsNestedCalleeOnly(t *testing.T) {
	// g()() : outer callee is a call expr (no name); inner is g.
	body := parseFuncBody(t, `package p
func f() { g()() }
`, "f")
	got := collectCalledIdents(body)
	if !hasAll(got, "g") {
		t.Errorf("collectCalledIdents = %v, want to include g", got)
	}
}
