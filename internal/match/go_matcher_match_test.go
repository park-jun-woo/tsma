package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoMatcherMatchFound(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "pkg", "handler")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.go"), []byte("package handler\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler_test.go"), []byte("package handler\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &GoMatcher{}
	testFile, found := m.Match(dir, "pkg/handler/handler.go")
	if !found {
		t.Fatal("expected to find test file")
	}
	want := filepath.Join("pkg", "handler", "handler_test.go")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestGoMatcherMatchNotFound(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.go"), []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &GoMatcher{}
	_, found := m.Match(dir, "pkg/handler.go")
	if found {
		t.Error("expected no match when test file does not exist")
	}
}

func TestGoMatcherMatchNonGoFile(t *testing.T) {
	dir := t.TempDir()

	m := &GoMatcher{}
	_, found := m.Match(dir, "handler.py")
	if found {
		t.Error("expected no match for non-.go file")
	}
}

func TestGoMatcherMatchRootFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &GoMatcher{}
	testFile, found := m.Match(dir, "main.go")
	if !found {
		t.Fatal("expected to find test file at root")
	}
	if testFile != "main_test.go" {
		t.Errorf("testFile = %q, want %q", testFile, "main_test.go")
	}
}
