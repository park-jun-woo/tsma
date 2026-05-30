package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestFinishFirstPassIfDone_flipsAtEnd(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A"}, {Name: "B"}},
		CurrentIndex: 2,
	}
	finishFirstPassIfDone(sess)
	if !sess.FirstPassDone {
		t.Error("expected FirstPassDone=true when watermark past end")
	}
	if sess.CurrentIndex != 0 {
		t.Errorf("expected CurrentIndex reset to 0, got %d", sess.CurrentIndex)
	}
}

func TestFinishFirstPassIfDone_notYet(t *testing.T) {
	sess := &model.Session{
		Functions:    []model.Function{{Name: "A"}, {Name: "B"}},
		CurrentIndex: 1,
	}
	finishFirstPassIfDone(sess)
	if sess.FirstPassDone {
		t.Error("expected FirstPassDone=false mid-pass")
	}
	if sess.CurrentIndex != 1 {
		t.Errorf("expected CurrentIndex unchanged, got %d", sess.CurrentIndex)
	}
}

func TestFinishFirstPassIfDone_alreadyDoneNoReset(t *testing.T) {
	sess := &model.Session{
		Functions:     []model.Function{{Name: "A"}},
		CurrentIndex:  3, // a live interactive cursor that happens to be high
		FirstPassDone: true,
	}
	finishFirstPassIfDone(sess)
	if sess.CurrentIndex != 3 {
		t.Errorf("expected interactive cursor untouched, got %d", sess.CurrentIndex)
	}
}
