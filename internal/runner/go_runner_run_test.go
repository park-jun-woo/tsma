package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoRunnerRunNonexistentTestFile(t *testing.T) {
	r := &GoRunner{}
	_, err := r.Run("/tmp", "/nonexistent/path/handler_test.go")
	if err == nil {
		t.Fatal("expected error for nonexistent test file")
	}
}

func TestGoRunnerRunExtractsTestFuncs(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

import "testing"

func TestHello(t *testing.T) {}
`
	testPath := filepath.Join(dir, "hello_test.go")
	if err := os.WriteFile(testPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Verify ExtractTestFuncs works on the file (part of Run's logic)
	funcs, err := ExtractTestFuncs(testPath)
	if err != nil {
		t.Fatalf("ExtractTestFuncs: %v", err)
	}
	if len(funcs) != 1 || funcs[0] != "TestHello" {
		t.Errorf("funcs = %v, want [TestHello]", funcs)
	}

	// Verify buildGoTestArgs produces correct args
	pkgPath, err := resolveGoPkgPath(dir, testPath)
	if err != nil {
		t.Fatalf("resolveGoPkgPath: %v", err)
	}
	args := buildGoTestArgs(pkgPath, funcs)
	if len(args) < 4 {
		t.Fatalf("args too short: %v", args)
	}
	if args[0] != "test" {
		t.Errorf("args[0] = %q, want \"test\"", args[0])
	}
	if args[3] != "-run" {
		t.Errorf("args[3] = %q, want \"-run\"", args[3])
	}
	if args[4] != "TestHello" {
		t.Errorf("args[4] = %q, want \"TestHello\"", args[4])
	}
}
