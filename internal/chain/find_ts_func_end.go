//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Counts braces from a starting line to determine where a TS/JS function body ends
package chain

// findTSFuncEnd finds the end line of a function starting at lineIdx (0-based).
// Returns 1-indexed end line.
func findTSFuncEnd(lines []string, lineIdx int) int {
	depth := 0
	started := false

	for j := lineIdx; j < len(lines); j++ {
		delta := countBraces(lines[j])
		depth += delta
		if delta > 0 || started {
			started = true
		}
		if started && depth <= 0 {
			return j + 1 // 1-indexed
		}
	}

	return len(lines)
}
