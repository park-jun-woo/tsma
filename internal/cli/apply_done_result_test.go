package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestApplyDoneResult(t *testing.T) {
	fn := &model.Function{Name: "Foo", Status: model.StatusTodo}
	out := captureStdout(func() {
		applyDoneResult(fn, mkMatch("pkg/foo_test.go"),
			measureResult{outcome: outcomeDone, mtime: "m2", coveragePct: 92})
	})
	if fn.Status != model.StatusDone {
		t.Errorf("expected DONE, got %s", fn.Status)
	}
	if fn.CoveragePct != 92 {
		t.Errorf("expected coverage 92, got %.0f", fn.CoveragePct)
	}
	if !strings.Contains(out, "DONE") || !strings.Contains(out, "92%") {
		t.Errorf("expected DONE line with 92%%, got %q", out)
	}
}
