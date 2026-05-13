//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints uncovered branch locations with file and line number
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/coverage"
)

// printUncovered prints uncovered branch locations.
func printUncovered(branches []coverage.UncoveredBranch) {
	for _, ub := range branches {
		fmt.Printf("  UNCOVERED: %s:%d\n", ub.File, ub.Line)
	}
}
