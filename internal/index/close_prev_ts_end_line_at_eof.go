//ff:func feature=index type=helper control=sequence
//ff:what Closes the last TS function's EndLine at end of file
package index

import "github.com/park-jun-woo/tsma/internal/model"

// closePrevTSEndLineAtEOF sets the EndLine of the last function in the slice
// to lastNonEmptyLine if it still equals StartLine, used at end of file.
func closePrevTSEndLineAtEOF(functions []model.Function, lastNonEmptyLine int) {
	n := len(functions)
	if n == 0 {
		return
	}
	if functions[n-1].EndLine != functions[n-1].StartLine {
		return
	}
	functions[n-1].EndLine = lastNonEmptyLine
}
