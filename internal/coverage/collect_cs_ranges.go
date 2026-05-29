//ff:func feature=coverage type=helper control=sequence lang=csharp
//ff:what Extracts the C# function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectCsRanges extracts the function range from the function.
func collectCsRanges(fn *model.Function) []csFuncRange {
	if fn.File == "" {
		return nil
	}
	return []csFuncRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
