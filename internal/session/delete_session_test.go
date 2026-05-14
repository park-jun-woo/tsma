package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestDeleteExistingSession(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{Project: dir, Lang: "go"}
	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := Delete(dir); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	tsmaDir := filepath.Join(dir, ".tsma")
	if _, err := os.Stat(tsmaDir); !os.IsNotExist(err) {
		t.Error(".tsma directory still exists after Delete")
	}
}

func TestDeleteNonexistentSession(t *testing.T) {
	dir := t.TempDir()
	// No session exists — Delete should not fail.
	if err := Delete(dir); err != nil {
		t.Fatalf("Delete on nonexistent session: %v", err)
	}
}

func TestDeleteThenExistsReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{Project: dir, Lang: "go"}
	if err := Save(dir, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if !Exists(dir) {
		t.Fatal("Exists should be true before Delete")
	}

	if err := Delete(dir); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if Exists(dir) {
		t.Error("Exists returned true after Delete")
	}
}
