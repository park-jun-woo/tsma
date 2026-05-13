//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Determines the last line of a Django Python function or class by indentation
package endpoint

import "strings"

// findPyFuncEndDjango determines the last line of a Python function/class by indentation.
func findPyFuncEndDjango(lines []string, defIdx int, defIndent string) int {
	defIndentLen := effectiveIndent(defIndent)
	lastContentLine := defIdx + 1

	for i := defIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lineIndentLen := effectiveIndentStr(line)
		if lineIndentLen <= defIndentLen {
			break
		}

		lastContentLine = i + 1
	}

	return lastContentLine
}
