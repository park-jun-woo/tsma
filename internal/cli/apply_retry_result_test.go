package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestApplyRetryResult_keepsTodo(t *testing.T) {
	fn := &model.Function{Name: "Foo", Status: model.StatusTodo}
	out := captureStdout(func() {
		applyRetryResult(fn, mkMatch("pkg/foo_test.go"),
			measureResult{outcome: outcomeRetry, mtime: "m3", coveragePct: 50, attempt: 1})
	})
	// Partial must NOT be auto-promoted to DONE/PASS.
	if fn.Status != model.StatusTodo {
		t.Errorf("expected TODO retained, got %s", fn.Status)
	}
	if fn.Attempt != 1 || fn.CoveragePct != 50 {
		t.Errorf("expected attempt 1 / coverage 50, got %d / %.0f", fn.Attempt, fn.CoveragePct)
	}
	if !strings.Contains(out, "PARTIAL") || !strings.Contains(out, "attempt 1") {
		t.Errorf("expected PARTIAL line, got %q", out)
	}
}
