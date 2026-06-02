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
	funcs := make([]*model.Function, len(sess.Functions))
	for i := range sess.Functions {
		funcs[i] = &sess.Functions[i]
	}
	batchMeasureOtherFuncs(root, sess.Lang, funcs)
}

// batchMeasureOtherFuncs is the core single-function-loop measurement over an
// explicit slice, so it can measure the whole session or just the
// newly-reconciled subset (Phase012 rescan).
func batchMeasureOtherFuncs(root, lang string, funcs []*model.Function) {
	matcher := match.NewFuncMatcher(lang)
	checker := coverage.NewChecker(lang)
	total := len(funcs)

	for i, fn := range funcs {
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
