package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestAdvanceCursor_increments(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		CurrentIndex: 1,
	}
	advanceCursor(sess)
	if sess.CurrentIndex != 2 {
		t.Errorf("expected CurrentIndex=2, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceCursor_wraps(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A"}, {Name: "B"}},
		CurrentIndex: 1,
	}
	advanceCursor(sess)
	if sess.CurrentIndex != 0 {
		t.Errorf("expected wrap to 0, got %d", sess.CurrentIndex)
	}
}

func TestAdvanceCursor_emptyNoop(t *testing.T) {
	sess := &model.Session{Functions: []model.Function{}, CurrentIndex: 0}
	advanceCursor(sess)
	if sess.CurrentIndex != 0 {
		t.Errorf("expected 0 for empty session, got %d", sess.CurrentIndex)
	}
}
