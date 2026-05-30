package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestSurfaceNextInteractiveTodo_showsNext(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass},
			{Name: "B", Status: model.StatusTodo},
		},
		CurrentIndex: 0,
	}
	out := captureStdout(func() { surfaceNextInteractiveTodo(sess) })
	if !strings.Contains(out, "B") {
		t.Errorf("expected next TODO B surfaced, got %q", out)
	}
}

func TestSurfaceNextInteractiveTodo_allComplete(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A", Status: model.StatusPass}},
		CurrentIndex: 0,
	}
	out := captureStdout(func() { surfaceNextInteractiveTodo(sess) })
	if !strings.Contains(out, "All functions complete!") {
		t.Errorf("expected completion banner, got %q", out)
	}
}
