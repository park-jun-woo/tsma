//ff:func feature=cli type=helper control=selection
//ff:what First-pass watermark scan: measures each function once, advancing past partials/untested without blocking
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// runNextFirstPass drives the first pass after a reset: CurrentIndex is a
// forward-only watermark over 0..N. Each function is measured once with its
// existing tests. 100% becomes PASS; a partial (test present but <100%) or an
// untested function stays TODO but the watermark still advances, so a single
// partial never gates the functions behind it (BUG-002). When the watermark
// passes the end, the session flips to interactive mode.
func runNextFirstPass(root string, sess *model.Session) error {
	fn := advanceToNext(sess)
	if fn == nil {
		finishFirstPassIfDone(sess)
		printAllComplete()
		return saveSession(root, sess)
	}

	changed, tm := detectTestChange(root, sess.Lang, fn)
	testFile := representativeTestFile(tm)

	if len(tm.Files) == 0 {
		// No test attributed: TODO, surface a hint, advance the watermark.
		printTodoFunction(fn, "")
		if misnamed, canonical, found := match.FindMisnamedTest(root, sess.Lang, fn.File); found {
			printRenameInstruction(misnamed, canonical)
		} else {
			printNextInstruction()
		}
		sess.CurrentIndex++
		return finishFirstPassAndSave(root, sess)
	}

	if !changed {
		// Already measured at this mtime (rare in a fresh pass): keep TODO and
		// advance so the watermark does not stall.
		printTodoFunction(fn, testFile)
		printNextInstruction()
		sess.CurrentIndex++
		return finishFirstPassAndSave(root, sess)
	}

	result := runAndMeasure(root, sess.Lang, fn, tm)

	switch result.outcome {
	case outcomeTestFail:
		// Failing test: do NOT advance the watermark; the function must be
		// fixed before the first pass continues past it.
		applyTestFailResult(fn, tm, result, testFile)
		sess.RecalcSummary()
		return saveSession(root, sess)

	case outcomePass:
		applyPassResult(fn, tm, result)
		sess.CurrentIndex++

	case outcomeDone:
		applyDoneResult(fn, tm, result)
		sess.CurrentIndex++

	case outcomeRetry:
		// Partial: keep TODO (no auto-demotion to DONE) but advance the
		// watermark so the rest of the first pass is still measured.
		applyRetryResult(fn, tm, result)
		sess.CurrentIndex++
	}

	finishFirstPassIfDone(sess)
	if !sess.FirstPassDone {
		// advanceToNext moves the watermark over any PASS/DONE we just produced
		// to the next TODO. If it returns nil the watermark hit the end, so the
		// first pass is now complete.
		if next := advanceToNext(sess); next != nil {
			printContinueInstruction()
			printTodoFunction(next, next.TestFile)
			printNextInstruction()
		}
		finishFirstPassIfDone(sess)
	}

	sess.RecalcSummary()
	return saveSession(root, sess)
}
