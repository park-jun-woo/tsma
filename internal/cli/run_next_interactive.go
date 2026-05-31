//ff:func feature=cli type=helper control=selection
//ff:what Interactive rotating-cursor mode after the first pass: surface one TODO per call, never stall on a partial
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// runNextInteractive handles `tsma next` after the first pass. CurrentIndex is a
// rotating cursor over the remaining TODOs. A non-measuring TODO (untested, or an
// unchanged partial) is surfaced and the cursor advances (wrapping), so a single
// unchanged partial can never trap the TODOs behind it (BUG-002 ★). A TODO whose
// test changed is re-measured in place: PASS/DONE advances; a partial keeps the
// cursor so the just-edited function stays visible. When no TODO is measurable a
// remaining-TODO summary is appended so the user/agent can pick what to work on.
func runNextInteractive(root string, sess *model.Session) error {
	fn := advanceToNextTodo(sess)
	if fn == nil {
		printAllComplete()
		return saveSession(root, sess)
	}

	changed, tm := detectTestChange(root, sess.Lang, fn)
	testFile := representativeTestFile(tm)

	if len(tm.Files) == 0 {
		// Untested: surface a write/rename hint, then advance the cursor so the
		// next call shows the next TODO instead of re-pinning this one.
		printTodoFunction(fn, "")
		if misnamed, canonical, found := match.FindMisnamedTest(root, sess.Lang, fn.File); found {
			printRenameInstruction(misnamed, canonical)
		} else {
			printNextInstruction()
		}
		advanceCursor(sess)
		maybePrintNoProgressSummary(root, sess)
		sess.RecalcSummary()
		return saveSession(root, sess)
	}

	if !changed {
		// Unchanged partial: count this presentation toward the auto-DONE
		// threshold. At the threshold it is accepted as best-effort DONE;
		// otherwise it is re-surfaced and the cursor rotates so a single partial
		// never traps the TODOs behind it (BUG-002 ★).
		return resurfacePartial(root, sess, fn, tm, testFile)
	}

	result := runAndMeasure(root, sess.Lang, fn, tm, effectiveMaxAttempts(sess))

	switch result.outcome {
	case outcomeTestFail:
		applyTestFailResult(fn, tm, result, testFile)

	case outcomePass:
		applyPassResult(fn, tm, result)
		advanceCursor(sess)
		surfaceNextInteractiveTodo(sess)

	case outcomeDone:
		applyDoneResult(fn, tm, result)
		advanceCursor(sess)
		surfaceNextInteractiveTodo(sess)

	case outcomeRetry:
		// Partial after a real edit: keep TODO (no auto-DONE) and keep the
		// cursor on it so the user keeps seeing the function they just touched.
		applyRetryResult(fn, tm, result)
		printTodoFunction(fn, testFile)
		printNextInstruction()
	}

	sess.RecalcSummary()
	return saveSession(root, sess)
}
