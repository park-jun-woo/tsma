//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Filters functions to exclude dead code functions
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// filterNonDead returns only non-dead functions.
func filterNonDead(functions []model.Function) []model.Function {
	var result []model.Function
	for _, fn := range functions {
		if fn.Dead {
			continue
		}
		result = append(result, fn)
	}
	return result
}
