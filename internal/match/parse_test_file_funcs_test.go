package match

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

func writeTmpFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseTestFileFuncs(t *testing.T) {
	path := writeTmpFile(t, "x_test.go", `package p
import "testing"
func TestA(t *testing.T) {}
func helper() {}
type T struct{}
func (T) Method() {}
var x = 1
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatalf("parseTestFileFuncs: %v", err)
	}
	for _, want := range []string{"TestA", "helper", "Method"} {
		if _, ok := funcs[want]; !ok {
			t.Errorf("missing func %q in %v", want, keys(funcs))
		}
	}
	// var x is not a func decl.
	if _, ok := funcs["x"]; ok {
		t.Errorf("non-func %q should not be present", "x")
	}
}

func TestParseTestFileFuncsParseError(t *testing.T) {
	path := writeTmpFile(t, "bad_test.go", `package p
func TestA(t *testing.T) { this is not go
`)
	if _, err := parseTestFileFuncs(path); err == nil {
		t.Fatal("expected parse error for malformed source, got nil")
	}
}

func TestParseTestFileFuncsMissingFile(t *testing.T) {
	if _, err := parseTestFileFuncs(filepath.Join(t.TempDir(), "nope_test.go")); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseTestFileFuncsCollisionFirstWins(t *testing.T) {
	path := writeTmpFile(t, "coll_test.go", `package p
type A struct{}
type B struct{}
func (A) Do() { FromA() }
func (B) Do() { FromB() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatalf("parseTestFileFuncs: %v", err)
	}
	fd, ok := funcs["Do"]
	if !ok {
		t.Fatal("Do not found")
	}
	// First declaration (A.Do) wins, so it calls FromA.
	called := collectCalledIdents(fd.Body)
	if !hasAll(called, "FromA") {
		t.Errorf("first-wins: Do body callees = %v, want FromA", called)
	}
	if _, bad := called["FromB"]; bad {
		t.Errorf("first-wins violated: got FromB in %v", called)
	}
}

func keys(m map[string]*ast.FuncDecl) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
