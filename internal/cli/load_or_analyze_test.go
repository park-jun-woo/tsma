package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestLoadOrAnalyze_existingSession(t *testing.T) {
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

	loaded, fresh, err := loadOrAnalyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fresh {
		t.Error("expected fresh=false on existing-session load path")
	}
	if loaded.Lang != "go" {
		t.Errorf("expected lang=go, got %s", loaded.Lang)
	}
	if len(loaded.Functions) != 1 {
		t.Errorf("expected 1 function, got %d", len(loaded.Functions))
	}
}

func TestLoadOrAnalyze_corruptSession(t *testing.T) {
	dir := t.TempDir()

	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	// Invalid JSON: Exists is true but Load fails.
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{broken"), 0o644)

	_, _, err := loadOrAnalyze(dir)
	if err == nil {
		t.Fatal("expected error loading corrupt session")
	}
}

func TestLoadOrAnalyze_saveError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; read-only dir permission check does not apply")
	}
	dir := t.TempDir()

	// Valid Go project so analyzeProject succeeds, but make .tsma a file so
	// session.Save cannot create the directory/write the session.
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	// Create .tsma as a regular file so MkdirAll inside Save fails.
	os.WriteFile(filepath.Join(dir, ".tsma"), []byte("not a dir"), 0o644)

	_, _, err := loadOrAnalyze(dir)
	if err == nil {
		t.Fatal("expected error when session cannot be saved")
	}
}

func TestLoadOrAnalyze_noSession_emptyDir(t *testing.T) {
	dir := t.TempDir()

	// No session and no Go project -> should fail at analyzeProject
	_, _, err := loadOrAnalyze(dir)
	if err == nil {
		t.Fatal("expected error for empty dir with no session")
	}
}

func TestLoadOrAnalyze_noSession_goProject(t *testing.T) {
	dir := t.TempDir()

	// Create minimal Go project
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)

	sess, fresh, err := loadOrAnalyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess == nil {
		t.Fatal("expected non-nil session")
	}
	if !fresh {
		t.Error("expected fresh=true on analyze path")
	}
	if sess.Lang != "go" {
		t.Errorf("expected lang=go, got %s", sess.Lang)
	}

	// Verify session was saved
	sessPath := filepath.Join(dir, ".tsma", "session.json")
	if _, err := os.Stat(sessPath); err != nil {
		t.Error("expected session file to be saved on disk")
	}
}
