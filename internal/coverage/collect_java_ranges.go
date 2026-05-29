//ff:func feature=coverage type=helper control=sequence
//ff:what Extracts the Java function range from a single Function
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// collectJavaRanges extracts the function range from the function.
func collectJavaRanges(fn *model.Function) []javaFuncRange {
	if fn.File == "" {
		return nil
	}
	return []javaFuncRange{{
		file:      fn.File,
		startLine: fn.StartLine,
		endLine:   fn.EndLine,
		funcName:  fn.Name,
	}}
}
