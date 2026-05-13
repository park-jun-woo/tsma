//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Collects incomplete (non-dead TODO or PARTIAL) functions from the session
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// collectIncomplete returns non-dead functions with TODO or PARTIAL status.
func collectIncomplete(functions []model.Function) []model.Function {
	var result []model.Function
	for _, fn := range functions {
		if fn.Dead {
			continue
		}
		if fn.Status == model.StatusTodo || fn.Status == model.StatusPartial {
			result = append(result, fn)
		}
	}
	return result
}
