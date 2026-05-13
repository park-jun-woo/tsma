//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts Go function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

func collectRanges(fn *model.Function) []funcRange {
	if fn.File == "" {
		return nil
	}
	return []funcRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
