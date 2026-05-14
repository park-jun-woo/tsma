package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestRunNext_allComplete(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Pass: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runNext(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestRunNext_todoNoTestFile(t *testing.T) {
	dir := t.TempDir()

	// Create source file but no test file
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg\nfunc Foo() {}"), 0o644)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 2, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runNext(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}
