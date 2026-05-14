package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestRunList_noSession(t *testing.T) {
	dir := t.TempDir()

	// Change to temp dir so getProjectRoot returns it
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runList(nil, nil)
	if err == nil {
		t.Fatal("expected error when no session exists")
	}
}

func TestRunList_withSession(t *testing.T) {
	dir := t.TempDir()

	// Create a valid session file
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 2, Pass: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Reset page to valid value
	listPage = 1

	output := captureStdout(func() {
		err := runList(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestRunList_pageOutOfRange(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	listPage = 999

	output := captureStdout(func() {
		err := runList(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}
