//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts Rust function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectRsRanges extracts the function range from the function.
func collectRsRanges(fn *model.Function) []rsFuncRange {
	if fn.File == "" {
		return nil
	}
	return []rsFuncRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
