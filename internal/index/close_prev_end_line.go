//ff:func feature=index type=helper control=sequence
//ff:what Closes the previous function's EndLine if it is still zero
package index

import "github.com/park-jun-woo/tsma/internal/model"

// closePrevEndLine sets the EndLine of the last function in the slice
// if it is still zero, indicating it hasn't been closed yet.
func closePrevEndLine(functions []model.Function, lastNonEmptyLine int) {
	n := len(functions)
	if n == 0 {
		return
	}
	if functions[n-1].EndLine != 0 {
		return
	}
	functions[n-1].EndLine = lastNonEmptyLine
}
