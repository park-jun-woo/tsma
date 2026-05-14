package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestRunStatus_noSession(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runStatus(nil, nil)
	if err == nil {
		t.Fatal("expected error when no session exists")
	}
}

func TestRunStatus_withSession(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusDone, CoveragePct: 80},
			{Name: "C", Status: model.StatusTodo},
		},
		Summary: model.Summary{Total: 3, Pass: 1, Done: 1, Todo: 1},
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
		err := runStatus(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "3 functions") {
		t.Errorf("expected '3 functions' in output, got %q", output)
	}
	if !strings.Contains(output, "PASS") {
		t.Error("expected PASS in output")
	}
	if !strings.Contains(output, "DONE") {
		t.Error("expected DONE in output")
	}
	if !strings.Contains(output, "TODO") {
		t.Error("expected TODO in output")
	}
	if !strings.Contains(output, "Coverage average") {
		t.Error("expected coverage average in output")
	}
}

func TestRunStatus_emptyFunctions(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project:   dir,
		Lang:      "go",
		Functions: []model.Function{},
		Summary:   model.Summary{Total: 0},
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
		err := runStatus(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "No functions found") {
		t.Errorf("expected 'No functions found' in output, got %q", output)
	}
}
