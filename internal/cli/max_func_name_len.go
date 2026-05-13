//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Finds the longest function name length in a slice for alignment
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// maxFuncNameLen returns the length of the longest function name.
func maxFuncNameLen(functions []model.Function) int {
	maxLen := 0
	for _, fn := range functions {
		if len(fn.Name) > maxLen {
			maxLen = len(fn.Name)
		}
	}
	return maxLen
}
