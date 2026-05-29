package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestTestMtimeOf_existingFile verifies the modification time of an existing
// file is returned, matching the file's actual mtime.
func TestTestMtimeOf_existingFile(t *testing.T) {
	root := t.TempDir()
	rel := "pkg/foo_test.go"
	abs := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	want := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(abs, want, want); err != nil {
		t.Fatal(err)
	}

	got := testMtimeOf(root, rel)
	if !got.Equal(want) {
		t.Errorf("testMtimeOf = %v, want %v", got, want)
	}
}

// TestTestMtimeOf_missingFile verifies a missing file yields the zero time.
func TestTestMtimeOf_missingFile(t *testing.T) {
	root := t.TempDir()
	got := testMtimeOf(root, "does_not_exist_test.go")
	if !got.IsZero() {
		t.Errorf("testMtimeOf = %v, want zero time", got)
	}
}
