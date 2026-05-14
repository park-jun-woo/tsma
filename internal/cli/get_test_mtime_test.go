package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetTestMtime_existingFile(t *testing.T) {
	dir := t.TempDir()
	relPath := "pkg/foo_test.go"
	absPath := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absPath, []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := getTestMtime(dir, relPath)
	if result == "" {
		t.Fatal("expected non-empty mtime string")
	}
	// Verify it parses as RFC3339
	_, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Errorf("expected RFC3339 format, got %q: %v", result, err)
	}
}

func TestGetTestMtime_nonexistentFile(t *testing.T) {
	dir := t.TempDir()
	result := getTestMtime(dir, "nonexistent_test.go")
	if result != "" {
		t.Errorf("expected empty string for nonexistent file, got %q", result)
	}
}

func TestGetTestMtime_emptyRoot(t *testing.T) {
	result := getTestMtime("", "some_test.go")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
