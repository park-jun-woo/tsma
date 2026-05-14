package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Login", Name: "Login", Status: model.StatusPass},
		},
	}
	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Lang != "go" {
		t.Errorf("Lang = %q, want %q", loaded.Lang, "go")
	}
	if len(loaded.Functions) != 1 {
		t.Fatalf("Functions count = %d, want 1", len(loaded.Functions))
	}
	if loaded.Functions[0].QualifiedName != "pkg.Login" {
		t.Errorf("QualifiedName = %q, want %q", loaded.Functions[0].QualifiedName, "pkg.Login")
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should fail for nonexistent session")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	sessionDir := filepath.Join(dir, ".tsma")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "session.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load should fail for invalid JSON")
	}
}

func TestLoadEmptyJSON(t *testing.T) {
	dir := t.TempDir()
	sessionDir := filepath.Join(dir, ".tsma")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "session.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Lang != "" {
		t.Errorf("Lang = %q, want empty", loaded.Lang)
	}
	if len(loaded.Functions) != 0 {
		t.Errorf("Functions count = %d, want 0", len(loaded.Functions))
	}
}
