package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileBytesSuccess(t *testing.T) {
	dir := t.TempDir()
	content := "hello world"
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := readFileBytes(path)
	if err != nil {
		t.Fatalf("readFileBytes: %v", err)
	}
	if string(data) != content {
		t.Errorf("content = %q, want %q", string(data), content)
	}
}

func TestReadFileBytesNonexistentPath(t *testing.T) {
	_, err := readFileBytes("/nonexistent/path/file.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReadFileBytesEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := readFileBytes(path)
	if err != nil {
		t.Fatalf("readFileBytes: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty, got %d bytes", len(data))
	}
}
