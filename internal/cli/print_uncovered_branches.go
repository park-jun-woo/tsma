//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints uncovered branch lines
package cli

import "fmt"

// printUncoveredBranches prints each uncovered branch line.
func printUncoveredBranches(branches []int) {
	for _, line := range branches {
		fmt.Printf("  uncovered branch: line %d\n", line)
	}
}
