//ff:func feature=cli type=helper control=sequence
//ff:what Counts an unchanged partial's presentation and auto-DONEs it at the threshold, else re-surfaces it
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// resurfacePartial handles an unchanged (not re-measured) partial function: it
// counts this presentation toward the auto-DONE threshold. At the threshold it
// accepts the function as best-effort DONE, keeping its last measured coverage
// and test mtime (no re-measure). Below the threshold it re-surfaces the function
// as TODO and rotates the cursor so a single partial never traps the TODOs behind
// it (BUG-002). Only test-bearing partials reach here (untested functions are
// handled earlier), so this never auto-DONEs an untested function.
func resurfacePartial(root string, sess *model.Session, fn *model.Function, tm match.TestMatch, testFile string) error {
	fn.Attempt++
	if attemptOutcome(fn.Attempt, effectiveMaxAttempts(sess)) == outcomeDone {
		applyDoneResult(fn, tm, measureResult{
			outcome:     outcomeDone,
			coveragePct: fn.CoveragePct,
			attempt:     fn.Attempt,
			mtime:       fn.TestMtime,
		})
		advanceCursor(sess)
		surfaceNextInteractiveTodo(sess)
		sess.RecalcSummary()
		return saveSession(root, sess)
	}
	printTodoFunction(fn, testFile)
	printNextInstruction()
	advanceCursor(sess)
	maybePrintNoProgressSummary(root, sess)
	sess.RecalcSummary()
	return saveSession(root, sess)
}
