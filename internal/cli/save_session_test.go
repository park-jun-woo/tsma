package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSaveSession_ok(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".tsma", "tests"), 0o755)
	sess := &model.Session{Project: dir, Lang: "go"}
	if err := saveSession(dir, sess); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveSession_error(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".tsma"), 0o755)
	makeSaveFail(t, dir) // .tsma/tests is a regular file -> MkdirAll fails
	sess := &model.Session{Project: dir, Lang: "go"}
	if err := saveSession(dir, sess); err == nil {
		t.Fatal("expected save error")
	}
}
