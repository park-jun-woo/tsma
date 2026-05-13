//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints uncovered line numbers from coverage data
package cli

import "fmt"

// printUncoveredLines prints each uncovered branch line from coverage.
func printUncoveredLines(lines []int) {
	for _, line := range lines {
		fmt.Printf("    UNCOVERED: branch at line %d\n", line)
	}
}
