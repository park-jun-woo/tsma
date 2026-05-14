package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExistsTrue(t *testing.T) {
	dir := t.TempDir()
	sessionDir := filepath.Join(dir, ".tsma")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "session.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !Exists(dir) {
		t.Error("Exists returned false, want true")
	}
}

func TestExistsFalseNoDir(t *testing.T) {
	dir := t.TempDir()
	if Exists(dir) {
		t.Error("Exists returned true for empty dir, want false")
	}
}

func TestExistsFalseNoSessionFile(t *testing.T) {
	dir := t.TempDir()
	sessionDir := filepath.Join(dir, ".tsma")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if Exists(dir) {
		t.Error("Exists returned true when session.json missing, want false")
	}
}
