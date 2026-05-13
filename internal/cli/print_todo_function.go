//ff:func feature=cli type=helper control=sequence
//ff:what Prints a TODO function with file location and test file path
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printTodoFunction prints a TODO function with file and test file info.
func printTodoFunction(fn *model.Function, testFile string) {
	fmt.Printf("\n%s  TODO\n", fn.Name)
	fmt.Printf("  file: %s:%d-%d\n", fn.File, fn.StartLine, fn.EndLine)
	if testFile != "" {
		fmt.Printf("  test: %s\n", testFile)
	} else {
		fmt.Printf("  test: (not found)\n")
	}
}
