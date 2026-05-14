package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSaveCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Login", Status: model.StatusTodo},
		},
	}

	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Verify .tsma directory created
	tsmaDir := filepath.Join(dir, ".tsma")
	if _, err := os.Stat(tsmaDir); err != nil {
		t.Errorf(".tsma directory not created: %v", err)
	}

	// Verify tests subdirectory created
	testsDir := filepath.Join(tsmaDir, "tests")
	if _, err := os.Stat(testsDir); err != nil {
		t.Errorf("tests directory not created: %v", err)
	}

	// Verify session.json created
	sessionPath := filepath.Join(tsmaDir, "session.json")
	if _, err := os.Stat(sessionPath); err != nil {
		t.Errorf("session.json not created: %v", err)
	}
}

func TestSaveRecalcsSummary(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Login", Status: model.StatusPass},
			{Name: "Signup", Status: model.StatusDone},
			{Name: "Logout", Status: model.StatusTodo},
		},
		// Intentionally wrong summary to verify recalc
		Summary: model.Summary{Total: 0},
	}

	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// After save, summary should be recalculated
	if sess.Summary.Total != 3 {
		t.Errorf("Total = %d, want 3", sess.Summary.Total)
	}
	if sess.Summary.Pass != 1 {
		t.Errorf("Pass = %d, want 1", sess.Summary.Pass)
	}
	if sess.Summary.Done != 1 {
		t.Errorf("Done = %d, want 1", sess.Summary.Done)
	}
	if sess.Summary.Todo != 1 {
		t.Errorf("Todo = %d, want 1", sess.Summary.Todo)
	}
}

func TestSaveWritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project:   dir,
		Lang:      "python",
		CheckedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Functions: []model.Function{
			{QualifiedName: "mod.handler", Name: "handler", Status: model.StatusPass},
		},
	}

	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Load it back to verify roundtrip
	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if loaded.Lang != "python" {
		t.Errorf("Lang = %q, want %q", loaded.Lang, "python")
	}
	if len(loaded.Functions) != 1 {
		t.Errorf("Functions count = %d, want 1", len(loaded.Functions))
	}
}

func TestSaveOverwritesExisting(t *testing.T) {
	dir := t.TempDir()

	sess1 := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Login", Status: model.StatusTodo},
		},
	}
	if err := Save(dir, sess1); err != nil {
		t.Fatalf("Save first: %v", err)
	}

	sess2 := &model.Session{
		Project: dir,
		Lang:    "typescript",
		Functions: []model.Function{
			{Name: "Render", Status: model.StatusPass},
			{Name: "Update", Status: model.StatusDone},
		},
	}
	if err := Save(dir, sess2); err != nil {
		t.Fatalf("Save second: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Lang != "typescript" {
		t.Errorf("Lang = %q, want %q", loaded.Lang, "typescript")
	}
	if len(loaded.Functions) != 2 {
		t.Errorf("Functions count = %d, want 2", len(loaded.Functions))
	}
}
