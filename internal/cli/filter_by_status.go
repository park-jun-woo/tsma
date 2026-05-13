//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Filters functions by a specific status, excluding dead functions
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// filterByStatus returns non-dead functions matching the given status.
func filterByStatus(functions []model.Function, status string) []model.Function {
	var result []model.Function
	for _, fn := range functions {
		if fn.Dead || fn.Status != status {
			continue
		}
		result = append(result, fn)
	}
	return result
}
