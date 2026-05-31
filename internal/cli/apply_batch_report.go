//ff:func feature=cli type=helper control=sequence
//ff:what Sets a function's status from a batch coverage report (PASS at 100%, else measured TODO)
package cli

import (
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applyBatchReport sets a function's status from a batch-measured coverage
// report. 100% coverage becomes PASS; anything below stays TODO with the
// measured percentage. Either way Attempt is set to 1, recording that the
// first-scan batch measured the function once. The batch never auto-promotes a
// partial to DONE (that is reserved for later interactive re-measurement), so
// this is a strict PASS-or-TODO decision. The matched test files and their
// combined mtime are recorded so a subsequent `tsma next` does not re-detect an
// unchanged test as a fresh edit.
func applyBatchReport(root string, fn *model.Function, m match.TestMatch, report *coverage.Report) {
	setTestFiles(fn, m)
	fn.TestMtime = combinedTestMtime(root, m.Files)
	fn.FailOutput = ""
	fn.Attempt = 1
	if report.AllCovered {
		fn.Status = model.StatusPass
		fn.CoveragePct = 100
		return
	}
	fn.Status = model.StatusTodo
	fn.CoveragePct = report.TotalPct
}
