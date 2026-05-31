package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestResurfacePartial_convergesToDone verifies that repeatedly resurfacing an
// unchanged partial (no re-measure) counts each presentation and auto-accepts it
// as DONE once attempt reaches maxAttempts, keeping the last measured coverage.
func TestResurfacePartial_convergesToDone(t *testing.T) {
	root := t.TempDir()
	fn := &model.Function{
		Name:        "Foo",
		File:        "foo.go",
		Status:      model.StatusTodo,
		CoveragePct: 80.0,
		Attempt:     0,
		TestMtime:   "1234",
	}
	sess := &model.Session{
		Functions:   []model.Function{*fn},
		MaxAttempts: 3,
	}
	tm := mkMatch("foo_test.go")

	// attempt 1 -> still TODO
	if err := resurfacePartial(root, sess, fn, tm, "foo_test.go"); err != nil {
		t.Fatalf("resurface 1: %v", err)
	}
	if fn.Status != model.StatusTodo || fn.Attempt != 1 {
		t.Fatalf("after 1st: status=%v attempt=%d, want todo/1", fn.Status, fn.Attempt)
	}

	// attempt 2 -> still TODO
	if err := resurfacePartial(root, sess, fn, tm, "foo_test.go"); err != nil {
		t.Fatalf("resurface 2: %v", err)
	}
	if fn.Status != model.StatusTodo || fn.Attempt != 2 {
		t.Fatalf("after 2nd: status=%v attempt=%d, want todo/2", fn.Status, fn.Attempt)
	}

	// attempt 3 -> auto DONE, coverage/mtime preserved
	if err := resurfacePartial(root, sess, fn, tm, "foo_test.go"); err != nil {
		t.Fatalf("resurface 3: %v", err)
	}
	if fn.Status != model.StatusDone || fn.Attempt != 3 {
		t.Fatalf("after 3rd: status=%v attempt=%d, want done/3", fn.Status, fn.Attempt)
	}
	if fn.CoveragePct != 80.0 {
		t.Fatalf("coverage should be preserved: got %v, want 80.0", fn.CoveragePct)
	}
	if fn.TestMtime != "1234" {
		t.Fatalf("test mtime should be preserved: got %q, want 1234", fn.TestMtime)
	}
}

// TestResurfacePartial_maxOneImmediateDone verifies maxAttempts=1 accepts the
// function as DONE on its first resurface.
func TestResurfacePartial_maxOneImmediateDone(t *testing.T) {
	root := t.TempDir()
	fn := &model.Function{
		Name:        "Bar",
		File:        "bar.go",
		Status:      model.StatusTodo,
		CoveragePct: 50.0,
	}
	sess := &model.Session{
		Functions:   []model.Function{*fn},
		MaxAttempts: 1,
	}
	if err := resurfacePartial(root, sess, fn, mkMatch("bar_test.go"), "bar_test.go"); err != nil {
		t.Fatalf("resurface: %v", err)
	}
	if fn.Status != model.StatusDone || fn.Attempt != 1 {
		t.Fatalf("max=1 first resurface: status=%v attempt=%d, want done/1", fn.Status, fn.Attempt)
	}
}
