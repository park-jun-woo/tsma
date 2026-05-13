//ff:func feature=endpoint type=helper control=sequence
//ff:what Determines function body boundaries via brace counting from a starting line
package endpoint

// findExportedFuncBounds finds the start and end lines of a function starting at
// the given line index (0-based). Returns 1-indexed line numbers.
func findExportedFuncBounds(lines []string, startIdx int) (int, int) {
	if startIdx < 0 || startIdx >= len(lines) {
		return 1, 1
	}

	startLine := startIdx + 1
	endLine := findBraceEnd(lines, startIdx)
	if endLine == 0 {
		return startLine, len(lines)
	}
	return startLine, endLine
}
