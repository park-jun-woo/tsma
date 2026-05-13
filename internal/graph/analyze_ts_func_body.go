//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Scans lines of a TS/JS function body to extract and resolve call expressions
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// analyzeTSFuncBody scans the body of a TS/JS function for call expressions.
func analyzeTSFuncBody(lines []string, callerIdx int, relPath string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	startLine := functions[callerIdx].StartLine
	endLine := functions[callerIdx].EndLine
	if endLine <= startLine {
		endLine = findTSEndLine(lines, startLine-1)
	}

	for lineIdx := startLine; lineIdx < endLine && lineIdx <= len(lines); lineIdx++ {
		resolveTSCallsInLine(lines[lineIdx-1], callerIdx, relPath, imports, functions, idx)
	}
}
