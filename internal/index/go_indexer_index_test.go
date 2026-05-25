package index

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectGoFunctions(t *testing.T) {
	src := `package api

func HandleLogin() {
	return
}

func init() {}

func main() {}

func processRequest() {
	return
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "handler.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var functions []model.Function
	collectGoFunctions(f, fset, "internal/api/handler.go", "internal/api", &functions)

	if len(functions) != 2 {
		t.Fatalf("expected 2 functions (main/init excluded), got %d", len(functions))
	}

	names := map[string]bool{}
	for _, fn := range functions {
		names[fn.Name] = true
	}
	if !names["HandleLogin"] {
		t.Error("expected HandleLogin to be collected")
	}
	if !names["processRequest"] {
		t.Error("expected processRequest to be collected")
	}
	if names["main"] {
		t.Error("main should be excluded")
	}
	if names["init"] {
		t.Error("init should be excluded")
	}
}

func TestCollectGoFunctionsSkipsInterfaceMethod(t *testing.T) {
	src := `package api

type Handler interface {
	Login() error
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "handler.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	var functions []model.Function
	collectGoFunctions(f, fset, "handler.go", "", &functions)

	if len(functions) != 0 {
		t.Errorf("expected 0 functions for interface-only file, got %d", len(functions))
	}
}

func TestGoIndexerIndexWithTempGoFile(t *testing.T) {
	dir := t.TempDir()

	goContent := `package myapp

func Serve() {
	return
}

type Router struct{}

func (r *Router) AddRoute(path string) {
	return
}
`
	abs := filepath.Join(dir, "server.go")
	if err := os.WriteFile(abs, []byte(goContent), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(funcs))
	}

	var serveFound, addRouteFound bool
	for _, fn := range funcs {
		switch fn.Name {
		case "Serve":
			serveFound = true
			if !fn.Exported {
				t.Error("Serve should be exported")
			}
		case "AddRoute":
			addRouteFound = true
			if fn.QualifiedName != "Router.AddRoute" {
				t.Errorf("AddRoute QualifiedName = %q, want %q", fn.QualifiedName, "Router.AddRoute")
			}
		}
	}
	if !serveFound {
		t.Error("expected to find Serve")
	}
	if !addRouteFound {
		t.Error("expected to find AddRoute")
	}
}

func TestGoIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	// Create a Go file at root
	rootGo := `package myapp

func Keep() {
	return
}
`
	if err := os.WriteFile(filepath.Join(dir, "keep.go"), []byte(rootGo), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a Go file inside "ignored/" directory
	ignoredDir := filepath.Join(dir, "ignored")
	if err := os.MkdirAll(ignoredDir, 0o755); err != nil {
		t.Fatal(err)
	}
	ignoredGo := `package ignored

func Excluded() {
	return
}
`
	if err := os.WriteFile(filepath.Join(ignoredDir, "excluded.go"), []byte(ignoredGo), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create .tsmignore excluding "ignored/"
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("ignored/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (ignored/ excluded), got %d", len(funcs))
	}
	if funcs[0].Name != "Keep" {
		t.Errorf("expected Keep, got %s", funcs[0].Name)
	}
}
