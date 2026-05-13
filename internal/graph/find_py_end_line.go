//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Estimates the end of a Python function by tracking indentation level
package graph

import "strings"

// findPyEndLine estimates the end of a Python function starting at lineIdx (0-based).
func findPyEndLine(lines []string, lineIdx int) int {
	if lineIdx >= len(lines) {
		return len(lines)
	}

	defIndent := pyCalcIndent(lines[lineIdx])

	for i := lineIdx + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		if pyCalcIndent(lines[i]) <= defIndent {
			return i
		}
	}
	return len(lines)
}
