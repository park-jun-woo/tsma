package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExistsTrue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !fileExists(path) {
		t.Error("fileExists returned false for existing file")
	}
}

func TestFileExistsFalse(t *testing.T) {
	if fileExists("/nonexistent/path/file.txt") {
		t.Error("fileExists returned true for nonexistent file")
	}
}

func TestFileExistsDirectory(t *testing.T) {
	dir := t.TempDir()
	// os.Stat succeeds for directories too
	if !fileExists(dir) {
		t.Error("fileExists returned false for directory")
	}
}
