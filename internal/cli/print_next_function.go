//ff:func feature=cli type=helper control=sequence
//ff:what Prints a single function in next-command format with status-dependent details
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printNextFunction prints a function in the next-command format.
func printNextFunction(fn *model.Function) {
	label := strings.ToUpper(fn.Status)
	fmt.Printf("%s  %s\n", fn.Name, label)
	fmt.Printf("  file: %s:%d-%d\n", fn.File, fn.StartLine, fn.EndLine)

	if fn.TestFile != "" {
		fmt.Printf("  test: %s\n", fn.TestFile)
	} else {
		testGuess := guessTestFile(fn.File)
		fmt.Printf("  test: %s (not found)\n", testGuess)
	}

	if fn.Status == model.StatusFail && fn.FailOutput != "" {
		printFailDetail(fn.FailOutput)
	}
}
