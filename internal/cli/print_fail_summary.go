//ff:func feature=cli type=helper control=sequence
//ff:what Prints trimmed fail output to stderr for concise display
package cli

import (
	"fmt"
	"os"
)

// printFailSummary prints trimmed fail output to stderr.
func printFailSummary(output string) {
	trimmed := trimOutput(output)
	if trimmed != "" {
		fmt.Fprintf(os.Stderr, "    %s\n", trimmed)
	}
}
