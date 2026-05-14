package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestAdvanceToNext_firstTodo(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
			{Name: "B", Status: model.StatusTodo},
		},
		CurrentIndex: 0,
	}
	fn := advanceToNext(sess)
	if fn == nil {
		t.Fatal("expected non-nil function")
	}
	if fn.Name != "A" {
		t.Errorf("expected A, got %s", fn.Name)
	}
	if sess.CurrentIndex != 0 {
		t.Errorf("expected CurrentIndex=0, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceToNext_skipsDone(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusDone},
			{Name: "B", Status: model.StatusPass},
			{Name: "C", Status: model.StatusTodo},
		},
		CurrentIndex: 0,
	}
	fn := advanceToNext(sess)
	if fn == nil {
		t.Fatal("expected non-nil function")
	}
	if fn.Name != "C" {
		t.Errorf("expected C, got %s", fn.Name)
	}
	if sess.CurrentIndex != 2 {
		t.Errorf("expected CurrentIndex=2, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceToNext_allComplete(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass},
			{Name: "B", Status: model.StatusDone},
		},
		CurrentIndex: 0,
	}
	fn := advanceToNext(sess)
	if fn != nil {
		t.Errorf("expected nil, got %s", fn.Name)
	}
}

func TestAdvanceToNext_emptySession(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{},
		CurrentIndex: 0,
	}
	fn := advanceToNext(sess)
	if fn != nil {
		t.Errorf("expected nil for empty session, got %v", fn)
	}
}

func TestAdvanceToNext_startFromMiddle(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
			{Name: "B", Status: model.StatusDone},
			{Name: "C", Status: model.StatusTodo},
		},
		CurrentIndex: 1,
	}
	fn := advanceToNext(sess)
	if fn == nil {
		t.Fatal("expected non-nil function")
	}
	if fn.Name != "C" {
		t.Errorf("expected C, got %s", fn.Name)
	}
}

func TestAdvanceToNext_indexBeyondLength(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "A", Status: model.StatusTodo},
		},
		CurrentIndex: 5,
	}
	fn := advanceToNext(sess)
	if fn != nil {
		t.Errorf("expected nil when index beyond length, got %v", fn)
	}
}
