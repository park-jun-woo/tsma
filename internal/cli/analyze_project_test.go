package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeProject_emptyDir(t *testing.T) {
	dir := t.TempDir()

	// analyzeProject on an empty dir should fail at detect
	_, err := analyzeProject(dir)
	if err == nil {
		t.Fatal("expected error for empty project dir")
	}
}

func TestAnalyzeProject_goProject(t *testing.T) {
	dir := t.TempDir()

	// Create a minimal Go project structure
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)

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
}
