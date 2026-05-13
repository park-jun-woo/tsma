//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Determines the last line of a Python function by indentation analysis
package endpoint

import "strings"

// findPyFuncEnd determines the last line of a Python function by indentation.
// defIdx is 0-based index of the def line, funcIndent is the indentation of the def line.
func findPyFuncEnd(lines []string, defIdx int, funcIndent string) int {
	defIndent := effectiveIndent(funcIndent)
	lastContentLine := defIdx + 1

	for i := defIdx + 1; i < len(lines); i++ {
		line := lines[i]

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		lineIndent := effectiveIndentStr(line)
		if lineIndent <= defIndent {
			break
		}

		lastContentLine = i + 1
	}

	return lastContentLine
}
