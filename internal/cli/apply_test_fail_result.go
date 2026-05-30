//ff:func feature=cli type=helper control=sequence
//ff:what Records a failing-test outcome on the function and prints the FAIL diagnostics
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// applyTestFailResult records a failing test: it stores the failure output and
// test files (status stays TODO) and prints the FAIL diagnostics plus the TODO
// function. Shared by both next modes.
func applyTestFailResult(fn *model.Function, tm match.TestMatch, result measureResult, testFile string) {
	setTestFiles(fn, tm)
	fn.TestMtime = result.mtime
	fn.FailOutput = result.failOutput
	fmt.Fprintf(os.Stderr, "FAIL  %s\n", fn.Name)
	fmt.Fprintf(os.Stderr, "  %s\n", result.failOutput)
	printTodoFunction(fn, testFile)
	printNextInstruction()
}
