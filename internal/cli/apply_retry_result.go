//ff:func feature=cli type=helper control=sequence
//ff:what Records a PARTIAL (retry) outcome, keeping the function TODO, and prints the PARTIAL line
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applyRetryResult records a partial measurement: it stores coverage/attempt and
// keeps fn TODO (never auto-demoted to DONE), then prints the PARTIAL line and
// uncovered branches. Shared by both next modes.
func applyRetryResult(fn *model.Function, tm match.TestMatch, result measureResult) {
	setTestFiles(fn, tm)
	fn.TestMtime = result.mtime
	fn.Attempt = result.attempt
	fn.CoveragePct = result.coveragePct
	fn.FailOutput = ""
	fmt.Printf("PARTIAL  %s  %.0f%%  (attempt %d — improve coverage)\n",
		fn.Name, result.coveragePct, result.attempt)
	printUncovered(result.uncovered)
}
