//ff:type feature=cli type=model
//ff:what Holds the outcome of a test run and coverage measurement
package cli

import "github.com/park-jun-woo/tsma/internal/coverage"

const (
	outcomeTestFail = "test_fail"
	outcomePass     = "pass"
	outcomeDone     = "done"
	outcomeRetry    = "retry"
)

type measureResult struct {
	outcome     string
	mtime       string
	coveragePct float64
	attempt     int
	failOutput  string
	uncovered   []coverage.UncoveredBranch
}
