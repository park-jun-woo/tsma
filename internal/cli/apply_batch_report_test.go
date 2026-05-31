//ff:test feature=cli
package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

func TestApplyBatchReport_FullCoveragePasses(t *testing.T) {
	fn := &model.Function{Name: "F", Status: model.StatusTodo}
	m := match.TestMatch{Files: []string{"f_test.go"}}
	applyBatchReport(t.TempDir(), fn, m, &coverage.Report{AllCovered: true, TotalPct: 100})
	if fn.Status != model.StatusPass {
		t.Errorf("want PASS, got %s", fn.Status)
	}
	if fn.CoveragePct != 100 {
		t.Errorf("want 100%%, got %.1f", fn.CoveragePct)
	}
	if fn.Attempt != 1 {
		t.Errorf("want Attempt=1, got %d", fn.Attempt)
	}
}

func TestApplyBatchReport_PartialStaysTodoNoAutoDone(t *testing.T) {
	fn := &model.Function{Name: "F", Status: model.StatusTodo}
	m := match.TestMatch{Files: []string{"f_test.go"}}
	applyBatchReport(t.TempDir(), fn, m, &coverage.Report{AllCovered: false, TotalPct: 62.5})
	if fn.Status != model.StatusTodo {
		t.Errorf("partial must stay TODO (no auto-DONE), got %s", fn.Status)
	}
	if fn.CoveragePct != 62.5 {
		t.Errorf("want measured 62.5%%, got %.1f", fn.CoveragePct)
	}
	if fn.Attempt != 1 {
		t.Errorf("want Attempt=1, got %d", fn.Attempt)
	}
	if fn.TestFile != "f_test.go" {
		t.Errorf("want test file recorded, got %q", fn.TestFile)
	}
}
