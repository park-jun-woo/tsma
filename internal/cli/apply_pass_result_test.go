package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestApplyPassResult(t *testing.T) {
	fn := &model.Function{Name: "Foo", Status: model.StatusTodo}
	out := captureStdout(func() {
		applyPassResult(fn, mkMatch("pkg/foo_test.go"),
			measureResult{outcome: outcomePass, mtime: "m1", coveragePct: 100})
	})
	if fn.Status != model.StatusPass {
		t.Errorf("expected PASS, got %s", fn.Status)
	}
	if fn.CoveragePct != 100 || fn.TestMtime != "m1" {
		t.Errorf("expected coverage 100 / mtime m1, got %.0f / %q", fn.CoveragePct, fn.TestMtime)
	}
	if !strings.Contains(out, "PASS") || !strings.Contains(out, "Foo") {
		t.Errorf("expected PASS line, got %q", out)
	}
}
