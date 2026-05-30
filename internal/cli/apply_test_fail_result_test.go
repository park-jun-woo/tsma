package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestApplyTestFailResult(t *testing.T) {
	fn := &model.Function{Name: "Foo", Status: model.StatusTodo}
	out := captureStdoutAndStderr(func() {
		applyTestFailResult(fn, mkMatch("pkg/foo_test.go"),
			measureResult{outcome: outcomeTestFail, mtime: "m4", failOutput: "boom"},
			"pkg/foo_test.go")
	})
	if fn.FailOutput != "boom" || fn.TestMtime != "m4" {
		t.Errorf("expected failOutput/mtime recorded, got %q / %q", fn.FailOutput, fn.TestMtime)
	}
	if fn.Status != model.StatusTodo {
		t.Errorf("expected status to stay TODO, got %s", fn.Status)
	}
	if !strings.Contains(out, "FAIL") || !strings.Contains(out, "boom") {
		t.Errorf("expected FAIL diagnostics, got %q", out)
	}
}
