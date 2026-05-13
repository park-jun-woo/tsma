//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Returns the net brace depth change for a single line of code
package chain

// countBraces returns the net brace depth change for a line.
func countBraces(line string) int {
	delta := 0
	for _, ch := range line {
		if ch == '{' {
			delta++
		} else if ch == '}' {
			delta--
		}
	}
	return delta
}
