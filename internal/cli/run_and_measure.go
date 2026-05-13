//ff:func feature=cli type=helper control=sequence
//ff:what Runs tests and measures coverage returning a structured result
package cli

import (
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// runAndMeasure executes the test, then measures coverage if tests pass.
func runAndMeasure(root, lang string, fn *model.Function, testFile string) measureResult {
	r := runner.NewRunner(lang)
	mtime := getTestMtime(root, testFile)

	res, err := r.Run(root, testFile)
	if err != nil || !res.Pass {
		output := ""
		if err != nil {
			output = err.Error()
		} else {
			output = res.Output
		}
		return measureResult{outcome: outcomeTestFail, mtime: mtime, failOutput: output}
	}

	checker := coverage.NewChecker(lang)
	report, err := checker.Check(root, testFile, fn)
	if err != nil {
		return measureResult{outcome: outcomeTestFail, mtime: mtime, failOutput: err.Error()}
	}

	if report.AllCovered {
		return measureResult{outcome: outcomePass, mtime: mtime, coveragePct: 100}
	}

	attempt := fn.Attempt + 1
	if attempt >= 2 {
		return measureResult{
			outcome:     outcomeDone,
			mtime:       mtime,
			coveragePct: report.TotalPct,
			attempt:     attempt,
		}
	}

	return measureResult{
		outcome:     outcomeRetry,
		mtime:       mtime,
		coveragePct: report.TotalPct,
		attempt:     attempt,
		uncovered:   report.Uncovered,
	}
}
