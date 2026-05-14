package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyTestFileSuccess(t *testing.T) {
	dir := t.TempDir()

	// Create source test file
	srcContent := "package handler\n\nfunc TestLogin(t *testing.T) {}\n"
	srcPath := filepath.Join(dir, "handler_test.go")
	if err := os.WriteFile(srcPath, []byte(srcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	rel, err := CopyTestFile(dir, srcPath)
	if err != nil {
		t.Fatalf("CopyTestFile: %v", err)
	}

	expectedRel := filepath.Join(".tsma", "tests", "handler_test.go")
	if rel != expectedRel {
		t.Errorf("relative path = %q, want %q", rel, expectedRel)
	}

	// Verify content
	dst := filepath.Join(dir, rel)
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copy: %v", err)
	}
	if string(data) != srcContent {
		t.Errorf("content = %q, want %q", string(data), srcContent)
	}
}

func TestCopyTestFileCreatesDirectory(t *testing.T) {
	dir := t.TempDir()

	srcPath := filepath.Join(dir, "test_main.py")
	if err := os.WriteFile(srcPath, []byte("import unittest"), 0o644); err != nil {
		t.Fatal(err)
	}

	// .tsma/tests does not exist yet
	_, err := CopyTestFile(dir, srcPath)
	if err != nil {
		t.Fatalf("CopyTestFile: %v", err)
	}

	// Verify directory was created
	testsDir := filepath.Join(dir, ".tsma", "tests")
	if _, err := os.Stat(testsDir); err != nil {
		t.Errorf("tests directory not created: %v", err)
	}
}

func TestCopyTestFileNonexistentSource(t *testing.T) {
	dir := t.TempDir()
	_, err := CopyTestFile(dir, filepath.Join(dir, "nonexistent_test.go"))
	if err == nil {
		t.Fatal("CopyTestFile should fail for nonexistent source")
	}
}

func TestCopyTestFilePreservesContent(t *testing.T) {
	dir := t.TempDir()

	// Create a binary-like content
	srcContent := "line1\nline2\nline3\n\x00binary"
	srcPath := filepath.Join(dir, "data_test.go")
	if err := os.WriteFile(srcPath, []byte(srcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	rel, err := CopyTestFile(dir, srcPath)
	if err != nil {
		t.Fatalf("CopyTestFile: %v", err)
	}

	dst := filepath.Join(dir, rel)
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copy: %v", err)
	}
	if string(data) != srcContent {
		t.Error("file content was not preserved exactly")
	}
}
