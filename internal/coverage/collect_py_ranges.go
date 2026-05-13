//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts Python function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectPyRanges extracts function range from the function.
func collectPyRanges(fn *model.Function) []pyFuncRange {
	if fn.File == "" {
		return nil
	}
	return []pyFuncRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
