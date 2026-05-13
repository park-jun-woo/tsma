//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Counts opening and closing braces in a single line
package endpoint

// countBraces returns the number of '{' and '}' characters in a line.
func countBraces(line string) (open, close int) {
	for _, ch := range line {
		if ch == '{' {
			open++
		} else if ch == '}' {
			close++
		}
	}
	return
}
