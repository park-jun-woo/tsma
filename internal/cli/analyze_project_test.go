package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestAnalyzeProject_emptyDir(t *testing.T) {
	dir := t.TempDir()

	// analyzeProject on an empty dir should fail at detect
	_, err := analyzeProject(dir)
	if err == nil {
		t.Fatal("expected error for empty project dir")
	}
}

// TestAnalyzeProject_goNoTestFile exercises the match loop's not-found branch:
// a Go project with a function but no matching test file, so every function
// keeps an empty TestFile.
func TestAnalyzeProject_goNoTestFile(t *testing.T) {
	dir := t.TempDir()

	// Go project with a function but NO matching test file: exercises the
	// match loop's not-found branch (TestFile stays empty for all funcs).
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "lib.go"), []byte("package lib\n\nfunc Add(a, b int) int { return a + b }\n"), 0o644)

	sess, err := analyzeProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sess.Functions) == 0 {
		t.Fatal("expected at least one indexed function")
	}
	for _, fn := range sess.Functions {
		if fn.TestFile != "" {
			t.Errorf("function %s: expected empty TestFile, got %q", fn.Name, fn.TestFile)
		}
		if fn.Status != model.StatusTodo {
			t.Errorf("function %s: expected status %s, got %s", fn.Name, model.StatusTodo, fn.Status)
		}
	}
}

// TestAnalyzeProject_pythonProject exercises analyzeProject end-to-end for a
// non-Go language: detect must resolve python from requirements.txt and the
// python indexer must collect the defined function.
func TestAnalyzeProject_pythonProject(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "lib.py"), []byte("def add(a, b):\n    return a + b\n"), 0o644)

	sess, err := analyzeProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess.Lang != "python" {
		t.Errorf("expected lang=python, got %s", sess.Lang)
	}
	if len(sess.Functions) == 0 {
		t.Fatal("expected at least one indexed python function")
	}
	for _, fn := range sess.Functions {
		if fn.Status != model.StatusTodo {
			t.Errorf("function %s: expected status %s, got %s", fn.Name, model.StatusTodo, fn.Status)
		}
	}
}

func TestAnalyzeProject_goProject(t *testing.T) {
	dir := t.TempDir()

	// Create a minimal Go project structure with a function and a matching
	// test file so the match loop hits both the indexed-function path and the
	// found-test-file branch.
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n\nfunc Greet() string { return \"hi\" }\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n\nimport \"testing\"\n\nfunc TestGreet(t *testing.T) {}\n"), 0o644)

	sess, err := analyzeProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess == nil {
		t.Fatal("expected non-nil session")
	}
	if sess.Lang != "go" {
		t.Errorf("expected lang=go, got %s", sess.Lang)
	}
	if sess.Project != dir {
		t.Errorf("expected project=%s, got %s", dir, sess.Project)
	}
	if len(sess.Functions) == 0 {
		t.Fatal("expected at least one indexed function")
	}
	// Every function should have been marked TODO by the match loop.
	for _, fn := range sess.Functions {
		if fn.Status != model.StatusTodo {
			t.Errorf("function %s: expected status %s, got %s", fn.Name, model.StatusTodo, fn.Status)
		}
	}
	// At least one function should have its TestFile populated (found branch).
	foundTest := false
	for _, fn := range sess.Functions {
		if fn.TestFile != "" {
			foundTest = true
			break
		}
	}
	if !foundTest {
		t.Error("expected at least one function to have a matched TestFile")
	}
}
