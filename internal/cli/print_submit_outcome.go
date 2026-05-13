//ff:func feature=cli type=helper control=sequence
//ff:what Prints the DONE or PARTIAL outcome after a submit
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

// printSubmitOutcome prints the result message after a submit.
func printSubmitOutcome(fn *model.Function, funcName string, report *coverage.Report, sess *model.Session) {
	if fn.Status == model.StatusDone {
		remaining := sess.Summary.Todo + sess.Summary.Partial
		fmt.Printf("DONE — %s complete (%d remaining)\n", funcName, remaining)
		return
	}

	uncoveredCount := len(report.Uncovered)
	fmt.Printf("PARTIAL — %s needs %d more branch(es) covered", funcName, uncoveredCount)
	if fn.RetryCount > 0 {
		fmt.Printf(" (retry %d/2)", fn.RetryCount)
	}
	fmt.Println()
}
