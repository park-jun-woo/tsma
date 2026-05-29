package match

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// parseFuncDecl parses src (a full file) and returns the top-level func decl
// named name.
func parseFuncDecl(t *testing.T, src, name string) *ast.FuncDecl {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return fd
		}
	}
	t.Fatalf("func %q not found", name)
	return nil
}

func TestMergeHelperRefs(t *testing.T) {
	helper := parseFuncDecl(t, `package p
func helper() {
	ReadSource()
	Validate()
}
`, "helper")

	refs := map[string]struct{}{"Existing": {}}
	mergeHelperRefs(refs, helper)

	if !hasAll(refs, "Existing", "ReadSource", "Validate") {
		t.Errorf("mergeHelperRefs result = %v, want to include Existing,ReadSource,Validate", refs)
	}
	if len(refs) != 3 {
		t.Errorf("len = %d, want 3: %v", len(refs), refs)
	}
}

func TestMergeHelperRefsOverlap(t *testing.T) {
	helper := parseFuncDecl(t, `package p
func helper() { ReadSource() }
`, "helper")

	refs := map[string]struct{}{"ReadSource": {}}
	mergeHelperRefs(refs, helper)

	if len(refs) != 1 || !hasAll(refs, "ReadSource") {
		t.Errorf("mergeHelperRefs overlap = %v, want single ReadSource", refs)
	}
}
