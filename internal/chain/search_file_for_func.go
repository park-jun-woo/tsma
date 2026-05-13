//ff:func feature=chain type=implementation control=sequence
//ff:what Reads a file and regex-matches a function definition, returning location if found
package chain

import (
	"os"
	"regexp"
	"strings"
)

// searchFileForFunc searches a single file for a function definition.
func searchFileForFunc(absPath, relPath, funcName string, pattern *regexp.Regexp) *tsFuncDef {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil
	}

	content := string(data)
	loc := pattern.FindStringIndex(content)
	if loc == nil {
		return nil
	}

	lines := strings.Split(content, "\n")
	lineNum := strings.Count(content[:loc[0]], "\n")

	startLine := lineNum + 1 // 1-indexed
	endLine := findTSFuncEnd(lines, lineNum)

	return &tsFuncDef{
		name:      funcName,
		file:      relPath,
		startLine: startLine,
		endLine:   endLine,
	}
}
