//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Finds the closing brace of a TS/JS function by tracking brace depth
package graph

// findTSEndLine finds the closing brace of a function starting at lineIdx (0-based).
func findTSEndLine(lines []string, lineIdx int) int {
	depth := 0
	for i := lineIdx; i < len(lines); i++ {
		depth = updateBraceDepth(lines[i], depth)
		if depth == 0 && i > lineIdx {
			return i + 1
		}
	}
	return len(lines)
}
