//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Scans lines of a Python function body to extract and resolve call expressions
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// analyzePyFuncBody scans the body of a Python function for call expressions.
func analyzePyFuncBody(lines []string, callerIdx int, relPath string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	startLine := functions[callerIdx].StartLine
	endLine := functions[callerIdx].EndLine
	if endLine <= startLine {
		endLine = findPyEndLine(lines, startLine-1)
	}

	for lineIdx := startLine; lineIdx < endLine && lineIdx <= len(lines); lineIdx++ {
		resolvePyCallsInLine(lines[lineIdx-1], callerIdx, relPath, imports, functions, idx)
	}
}
