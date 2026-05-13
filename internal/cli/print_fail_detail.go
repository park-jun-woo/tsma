//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints fail output lines indented under an error heading
package cli

import (
	"fmt"
	"strings"
)

// printFailDetail prints fail output lines indented under an error heading.
func printFailDetail(failOutput string) {
	fmt.Println("  error:")
	for _, line := range strings.Split(failOutput, "\n") {
		if line != "" {
			fmt.Printf("    %s\n", line)
		}
	}
}
