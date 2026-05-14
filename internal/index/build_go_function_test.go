package index

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestBuildGoFunctionExported(t *testing.T) {
	src := `package main

func HandleLogin() {
	return
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "handler.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	fd := f.Decls[0].(*ast.FuncDecl)
	fn := buildGoFunction(fd, fset, "internal/api/handler.go", "internal/api")

	if fn.QualifiedName != "internal/api.HandleLogin" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "internal/api.HandleLogin")
	}
	if fn.Name != "HandleLogin" {
		t.Errorf("Name = %q, want %q", fn.Name, "HandleLogin")
	}
	if fn.File != "internal/api/handler.go" {
		t.Errorf("File = %q, want %q", fn.File, "internal/api/handler.go")
	}
	if !fn.Exported {
		t.Error("expected Exported=true for HandleLogin")
	}
	if fn.StartLine != 3 {
		t.Errorf("StartLine = %d, want 3", fn.StartLine)
	}
	if fn.Status != "todo" {
		t.Errorf("Status = %q, want %q", fn.Status, "todo")
	}
}

func TestBuildGoFunctionUnexported(t *testing.T) {
	src := `package main

func helperFunc() {
	return
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "helper.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	fd := f.Decls[0].(*ast.FuncDecl)
	fn := buildGoFunction(fd, fset, "helper.go", "")

	if fn.Exported {
		t.Error("expected Exported=false for helperFunc")
	}
	if fn.QualifiedName != "helperFunc" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "helperFunc")
	}
}

func TestBuildGoFunctionWithReceiver(t *testing.T) {
	src := `package main

type Server struct{}

func (s *Server) Start() {
	return
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "server.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Find the FuncDecl (skip the GenDecl for the type)
	var fd *ast.FuncDecl
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			fd = fn
			break
		}
	}
	if fd == nil {
		t.Fatal("no FuncDecl found")
	}

	fn := buildGoFunction(fd, fset, "pkg/server.go", "pkg")

	if fn.QualifiedName != "pkg.Server.Start" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "pkg.Server.Start")
	}
	if !fn.Exported {
		t.Error("expected Exported=true for Start")
	}
}
