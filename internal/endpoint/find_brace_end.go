//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Counts braces from a starting line to find the closing brace line number
package endpoint

// findBraceEnd counts braces from startIdx to find where the block closes.
// Returns 1-indexed line number, or 0 if braces never balanced.
func findBraceEnd(lines []string, startIdx int) int {
	depth := 0
	started := false

	for j := startIdx; j < len(lines); j++ {
		open, close := countBraces(lines[j])
		depth += open - close
		if open > 0 {
			started = true
		}
		if started && depth <= 0 {
			return j + 1
		}
	}

	return 0
}
