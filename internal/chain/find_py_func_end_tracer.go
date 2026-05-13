//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Determines the last line of a Python function by comparing indentation levels
package chain

import "strings"

// findPyFuncEndTracer determines the last line of a Python function by indentation.
func findPyFuncEndTracer(lines []string, defIdx int, defIndent string) int {
	defIndentLen := pyEffectiveIndent(defIndent)
	lastContentLine := defIdx + 1 // 1-based

	for i := defIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lineIndentLen := pyEffectiveIndent(pyLeadingWhitespace(line))
		if lineIndentLen <= defIndentLen {
			break
		}

		lastContentLine = i + 1
	}

	return lastContentLine
}
