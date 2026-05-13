//ff:func feature=cli type=helper control=sequence
//ff:what Prints detailed status information for a single function
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printFuncDetail prints detailed information for a specific function.
func printFuncDetail(fn *model.Function) {
	status := strings.ToUpper(fn.Status)
	if fn.Dead {
		status = "DEAD"
	}
	fmt.Printf("%s (%s:%d-%d) — %s\n",
		fn.Name, fn.File, fn.StartLine, fn.EndLine, status)

	if fn.TestFile != "" {
		fmt.Printf("  test file: %s\n", fn.TestFile)
	}

	if fn.CoveragePct > 0 {
		fmt.Printf("  coverage: %.1f%%\n", fn.CoveragePct)
	}

	printUncoveredBranches(fn.UncoveredBranches)

	fmt.Printf("  callers: %d\n", len(fn.Callers))
	fmt.Printf("  callees: %d\n", len(fn.Callees))
}
