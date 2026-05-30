//ff:func feature=cli type=helper control=sequence
//ff:what Records a PASS outcome on the function and prints the PASS line
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applyPassResult marks fn as PASS at 100% coverage, records its test files and
// mtime, and prints the PASS line. Shared by both next modes.
func applyPassResult(fn *model.Function, tm match.TestMatch, result measureResult) {
	fn.Status = model.StatusPass
	setTestFiles(fn, tm)
	fn.TestMtime = result.mtime
	fn.CoveragePct = 100
	fn.FailOutput = ""
	fmt.Printf("PASS  %s  100%%\n", fn.Name)
}
