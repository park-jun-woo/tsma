//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Updates function status, coverage, and retry count based on the coverage report
package cli

import (
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applySubmitResult updates the function based on the coverage report.
func applySubmitResult(fn *model.Function, report *coverage.Report, prevPct float64) {
	fn.CoveragePct = report.TotalPct

	if report.AllCovered {
		fn.Status = model.StatusDone
		fn.UncoveredBranches = nil
		fn.RetryCount = 0
		return
	}

	fn.Status = model.StatusPartial
	fn.UncoveredBranches = nil
	for _, ub := range report.Uncovered {
		fn.UncoveredBranches = append(fn.UncoveredBranches, ub.Line)
	}

	if report.TotalPct > prevPct {
		fn.RetryCount = 0
	} else {
		fn.RetryCount++
	}
}
