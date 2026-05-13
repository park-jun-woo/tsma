//ff:func feature=cli type=helper control=sequence
//ff:what Trims test output to the first 5 lines for concise display
package cli

import "strings"

// trimOutput returns the first 5 lines of output for concise display.
func trimOutput(output string) string {
	const maxLines = 5
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > maxLines {
		lines = append(lines[:maxLines], "...")
	}
	return strings.Join(lines, "\n    ")
}
