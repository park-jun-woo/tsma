//ff:func feature=index type=helper control=sequence
//ff:what Closes the previous TS function's EndLine if it equals StartLine
package index

import "github.com/park-jun-woo/tsma/internal/model"

// closePrevTSEndLine sets the EndLine of the last function in the slice
// if it still equals StartLine, indicating it hasn't been closed yet.
func closePrevTSEndLine(functions []model.Function, lineNum, lastNonEmpty int) {
	n := len(functions)
	if n == 0 {
		return
	}
	if functions[n-1].EndLine != functions[n-1].StartLine {
		return
	}
	functions[n-1].EndLine = lastNonEmptyBeforeLine(lineNum, lastNonEmpty)
}
