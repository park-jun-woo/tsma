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

func TestGoIndexerIndexIgnoresFileAndSkipsParseErrors(t *testing.T) {
	dir := t.TempDir()

	// A valid Go file that should be indexed.
	good := `package myapp

func Keep() {
	return
}
`
	if err := os.WriteFile(filepath.Join(dir, "keep.go"), []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}

	// A specific ignored FILE (not a directory) -> exercises the file-ignore
	// branch that returns nil without SkipDir.
	ignoredFile := `package myapp

func Skipped() {
	return
}
`
	if err := os.WriteFile(filepath.Join(dir, "ignoreme.go"), []byte(ignoredFile), 0o644); err != nil {
		t.Fatal(err)
	}

	// A Go file with a syntax error -> parser.ParseFile fails -> skipped.
	broken := "package myapp\n\nfunc Broken( {\n"
	if err := os.WriteFile(filepath.Join(dir, "broken.go"), []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("ignoreme.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	names := map[string]bool{}
	for _, fn := range funcs {
		names[fn.Name] = true
	}
	if !names["Keep"] {
		t.Error("expected Keep to be indexed")
	}
	if names["Skipped"] {
		t.Error("Skipped should have been ignored via .tsmignore")
	}
	if names["Broken"] {
		t.Error("Broken should have been skipped due to parse error")
	}
}

func TestGoIndexerIndexWalkError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; unreadable-dir walk error does not apply")
	}
	dir := t.TempDir()

	good := "package myapp\n\nfunc Keep() {\n\treturn\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "keep.go"), []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}

	// An unreadable subdirectory causes filepath.Walk to invoke the walk fn with
	// a non-nil err for its entries, exercising the err != nil branch (which is
	// swallowed so indexing still succeeds).
	sub := filepath.Join(dir, "locked")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "inner.go"), []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("expected indexing to succeed despite walk error, got: %v", err)
	}
	found := false
	for _, fn := range funcs {
		if fn.Name == "Keep" {
			found = true
		}
	}
	if !found {
		t.Error("expected Keep to still be indexed")
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
