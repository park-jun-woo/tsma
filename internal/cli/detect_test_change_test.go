package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestDetectTestChange_noTestFile(t *testing.T) {
	dir := t.TempDir()
	// Create a source file but no test file
	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{
		Name: "Foo",
		File: "pkg/foo.go",
	}
	changed, testFile := detectTestChange(dir, "go", fn)
	if changed {
		t.Error("expected changed=false when no test file")
	}
	if testFile != "" {
		t.Errorf("expected empty testFile, got %q", testFile)
	}
}

func TestDetectTestChange_testFileExists_sameMtime(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}
	testPath := filepath.Join(srcDir, "foo_test.go")
	if err := os.WriteFile(testPath, []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Get the mtime and set it on the function
	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatal(err)
	}
	mtime := info.ModTime().Format(time.RFC3339)

	fn := &model.Function{
		Name:      "Foo",
		File:      "pkg/foo.go",
		TestMtime: mtime,
	}
	changed, testFile := detectTestChange(dir, "go", fn)
	if changed {
		t.Error("expected changed=false when mtime matches")
	}
	if testFile == "" {
		t.Error("expected non-empty testFile")
	}
}

func TestDetectTestChange_testFileExists_differentMtime(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo_test.go"), []byte("package pkg"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{
		Name:      "Foo",
		File:      "pkg/foo.go",
		TestMtime: "2000-01-01T00:00:00Z", // old mtime
	}
	changed, testFile := detectTestChange(dir, "go", fn)
	if !changed {
		t.Error("expected changed=true when mtime differs")
	}
	if testFile == "" {
		t.Error("expected non-empty testFile")
	}
}
