//ff:func feature=chain type=helper control=sequence
//ff:what Creates a pyFuncInfo from a regex match at the given line index
package chain

// buildPyFuncInfo creates a pyFuncInfo from a regex match at the given line index.
func buildPyFuncInfo(m []string, lines []string, lineIdx int, relPath string) *pyFuncInfo {
	indent := m[1]
	funcName := m[2]
	startLine := lineIdx + 1
	endLine := findPyFuncEndTracer(lines, lineIdx, indent)

	bodyEnd := endLine - 1
	if bodyEnd >= len(lines) {
		bodyEnd = len(lines) - 1
	}
	var bodyLines []string
	if lineIdx+1 <= bodyEnd {
		bodyLines = lines[lineIdx+1 : bodyEnd+1]
	}

	return &pyFuncInfo{
		name:      funcName,
		file:      relPath,
		startLine: startLine,
		endLine:   endLine,
		bodyLines: bodyLines,
	}
}
