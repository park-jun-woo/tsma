package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestAdvanceToNextTodo_atCursor(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass},
			{Name: "B", Status: model.StatusTodo},
			{Name: "C", Status: model.StatusTodo},
		},
		CurrentIndex: 1,
	}
	fn := advanceToNextTodo(sess)
	if fn == nil || fn.Name != "B" {
		t.Fatalf("expected B, got %v", fn)
	}
	if sess.CurrentIndex != 1 {
		t.Errorf("expected CurrentIndex=1, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceToNextTodo_wrapsAround(t *testing.T) {
	// Cursor at end, only TODO is at index 0 -> must wrap.
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
			{Name: "B", Status: model.StatusPass},
			{Name: "C", Status: model.StatusPass},
		},
		CurrentIndex: 2,
	}
	fn := advanceToNextTodo(sess)
	if fn == nil || fn.Name != "A" {
		t.Fatalf("expected A via wrap, got %v", fn)
	}
	if sess.CurrentIndex != 0 {
		t.Errorf("expected CurrentIndex=0 after wrap, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceToNextTodo_noTodoReturnsNil(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass},
			{Name: "B", Status: model.StatusDone},
		},
		CurrentIndex: 0,
	}
	if fn := advanceToNextTodo(sess); fn != nil {
		t.Errorf("expected nil, got %v", fn)
	}
}

func TestAdvanceToNextTodo_emptyReturnsNil(t *testing.T) {
	sess := &model.Session{Functions: []model.Function{}, CurrentIndex: 0}
	if fn := advanceToNextTodo(sess); fn != nil {
		t.Errorf("expected nil for empty session, got %v", fn)
	}
}

func TestAdvanceToNextTodo_outOfRangeCursorResets(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A", Status: model.StatusTodo}},
		CurrentIndex: 9,
	}
	fn := advanceToNextTodo(sess)
	if fn == nil || fn.Name != "A" {
		t.Fatalf("expected A after cursor reset, got %v", fn)
	}
	if sess.CurrentIndex != 0 {
		t.Errorf("expected CurrentIndex reset to 0, got %d", sess.CurrentIndex)
	}
}
