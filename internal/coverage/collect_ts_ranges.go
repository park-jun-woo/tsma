//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts TS function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectTSRanges extracts function range from the function.
func collectTSRanges(fn *model.Function) []tsFuncRange {
	if fn.File == "" {
		return nil
	}
	return []tsFuncRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
