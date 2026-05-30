package match

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestRecordSourceDecls(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "x.go", `package p
type GoFile struct{}
func (g *GoFile) M() {}
func Free() {}
var x = 1
`, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	r := &PkgSourceReceivers{byName: make(map[string]map[string]struct{})}
	recordSourceDecls(r, f)

	if _, ok := r.byName["M"]["GoFile"]; !ok {
		t.Errorf("method M should be recorded under receiver GoFile: %v", r.byName["M"])
	}
	if _, ok := r.byName["Free"][""]; !ok {
		t.Errorf("free func Free should be recorded under empty receiver: %v", r.byName["Free"])
	}
	// Non-func decls (the var) must not appear.
	if _, ok := r.byName["x"]; ok {
		t.Errorf("var decl must not be recorded")
	}
}
