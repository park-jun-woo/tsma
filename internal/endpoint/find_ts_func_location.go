//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Finds start and end lines of a named function in a TS/JS file by brace matching
package endpoint

import "regexp"

// findTSFuncLocation searches lines for a function/const definition matching the
// given name and returns its start and end lines (1-indexed). End line is estimated
// by counting braces from the start of the definition.
func findTSFuncLocation(lines []string, funcName string) (int, int) {
	funcPattern := regexp.MustCompile(
		`(?:export\s+)?(?:async\s+)?function\s+` + regexp.QuoteMeta(funcName) + `\b` +
			`|(?:export\s+)?(?:const|let|var)\s+` + regexp.QuoteMeta(funcName) + `\s*=`,
	)

	for i, line := range lines {
		if !funcPattern.MatchString(line) {
			continue
		}

		startLine := i + 1
		endLine := findBraceEnd(lines, i)
		if endLine == 0 {
			return startLine, len(lines)
		}
		return startLine, endLine
	}

	return 0, 0
}
