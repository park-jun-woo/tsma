package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoRunnerRunNonexistentTestFile(t *testing.T) {
	r := &GoRunner{}
	_, err := r.Run("/tmp", "/nonexistent/path/handler_test.go")
	if err == nil {
		t.Fatal("expected error for nonexistent test file")
	}
}

// TestGoRunnerRun_absError covers the filepath.Abs error branch (line 14):
// Abs fails only when os.Getwd() fails, forced by removing the cwd.
func TestGoRunnerRun_absError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	r := &GoRunner{}
	if _, err := r.Run("/proj", "rel_test.go"); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}

// TestGoRunnerRun_pkgPathError covers the resolveGoPkgPath error branch (line
// 19): a relative projectRoot cannot be made relative to the absolute test
// path.
func TestGoRunnerRun_pkgPathError(t *testing.T) {
	r := &GoRunner{}
	_, err := r.Run("relative-root", "some_test.go")
	if err == nil {
		t.Fatal("expected error when projectRoot is relative")
	}
	if !strings.Contains(err.Error(), "relative package") {
		t.Errorf("expected relative-package error, got: %v", err)
	}
}

func TestGoRunnerRun_passing(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/r\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "foo.go"), []byte("package r\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "foo_test.go"),
		[]byte("package r\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo() != 1 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	r := &GoRunner{}
	res, err := r.Run(dir, "foo_test.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Pass {
		t.Errorf("expected Pass=true, output: %s", res.Output)
	}
}

func TestGoRunnerRun_failing(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/r\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "foo.go"), []byte("package r\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "foo_test.go"),
		[]byte("package r\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tt.Fatal(\"boom\")\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	r := &GoRunner{}
	res, err := r.Run(dir, "foo_test.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Pass {
		t.Error("expected Pass=false for failing test")
	}
	if res.Output == "" {
		t.Error("expected non-empty output for failing test")
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
