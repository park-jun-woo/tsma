package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestFinishFirstPassAndSave_flipsAndRecounts(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".tsma", "tests"), 0o755)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass},
			{Name: "B", Status: model.StatusTodo},
		},
		CurrentIndex: 2, // watermark past end
	}
	if err := finishFirstPassAndSave(dir, sess); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sess.FirstPassDone {
		t.Error("expected FirstPassDone=true")
	}
	if sess.CurrentIndex != 0 {
		t.Errorf("expected cursor reset to 0, got %d", sess.CurrentIndex)
	}
	if sess.Summary.Pass != 1 || sess.Summary.Todo != 1 {
		t.Errorf("expected summary pass=1 todo=1, got %+v", sess.Summary)
	}
}

func TestFinishFirstPassAndSave_saveError(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".tsma"), 0o755)
	makeSaveFail(t, dir)
	sess := &model.Session{Project: dir, Lang: "go", CurrentIndex: 1}
	if err := finishFirstPassAndSave(dir, sess); err == nil {
		t.Fatal("expected save error")
	}
}
