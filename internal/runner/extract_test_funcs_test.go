package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTestFuncsMultiple(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

import "testing"

func TestLogin(t *testing.T) {}
func TestSignup(t *testing.T) {}
func helperFunc() {}
`
	path := filepath.Join(dir, "handler_test.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := ExtractTestFuncs(path)
	if err != nil {
		t.Fatalf("ExtractTestFuncs: %v", err)
	}
	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %v", len(funcs), funcs)
	}

	expected := map[string]bool{"TestLogin": true, "TestSignup": true}
	for _, f := range funcs {
		if !expected[f] {
			t.Errorf("unexpected func: %q", f)
		}
	}
}

func TestExtractTestFuncsNoTestFuncs(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

func helperFunc() {}
func BenchmarkLogin(b *testing.B) {}
`
	path := filepath.Join(dir, "helper.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := ExtractTestFuncs(path)
	if err != nil {
		t.Fatalf("ExtractTestFuncs: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 funcs, got %v", funcs)
	}
}

func TestExtractTestFuncsFileNotFound(t *testing.T) {
	_, err := ExtractTestFuncs("/nonexistent/path/test.go")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestExtractTestFuncsSingleFunc(t *testing.T) {
	dir := t.TempDir()
	content := `package foo

func TestOnly(t *testing.T) {
	// single test
}
`
	path := filepath.Join(dir, "single_test.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := ExtractTestFuncs(path)
	if err != nil {
		t.Fatalf("ExtractTestFuncs: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1", len(funcs))
	}
	if funcs[0] != "TestOnly" {
		t.Errorf("func = %q, want %q", funcs[0], "TestOnly")
	}
}

func TestExtractTestFuncsEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty_test.go")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := ExtractTestFuncs(path)
	if err != nil {
		t.Fatalf("ExtractTestFuncs: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 funcs for empty file, got %v", funcs)
	}
}
