//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Batch-measures non-Go functions one at a time via the single-function Check
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// batchMeasureOther measures every non-Go function at analysis time by looping
// the single-function coverage Check (no package-level optimization yet). The UX
// is identical to Go's batch: 100% becomes PASS, a matched partial becomes a
// measured TODO (Attempt=1), untested stays TODO. A function whose tests fail to
// run is left TODO and reported, without aborting the scan.
func batchMeasureOther(root string, sess *model.Session) {
	matcher := match.NewFuncMatcher(sess.Lang)
	checker := coverage.NewChecker(sess.Lang)
	total := len(sess.Functions)

	for i := range sess.Functions {
		fn := &sess.Functions[i]
		m, found := matcher.MatchFunc(root, fn)
		if !found || len(m.Files) == 0 {
			continue // untested, stays TODO
		}
		fmt.Fprintf(os.Stderr, "Measuring coverage: %d/%d (%s)\n", i+1, total, fn.Name)
		report, err := checker.Check(root, m, fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  skipped (%s): %s\n", fn.Name, firstLine(err.Error()))
			continue
		}
		applyBatchReport(root, fn, m, report)
	}
}
