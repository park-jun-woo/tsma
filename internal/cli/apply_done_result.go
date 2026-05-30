//ff:func feature=cli type=helper control=sequence
//ff:what Records a DONE outcome (best-effort accept after an improvement attempt) and prints the DONE line
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applyDoneResult marks fn as DONE at its measured coverage, records its test
// files and mtime, and prints the DONE line. This is the explicit best-effort
// accept (attempt>=2 after the user edited the test), not an auto-demotion.
func applyDoneResult(fn *model.Function, tm match.TestMatch, result measureResult) {
	fn.Status = model.StatusDone
	setTestFiles(fn, tm)
	fn.TestMtime = result.mtime
	fn.CoveragePct = result.coveragePct
	fn.FailOutput = ""
	fmt.Printf("DONE  %s  %.0f%%\n", fn.Name, result.coveragePct)
}
