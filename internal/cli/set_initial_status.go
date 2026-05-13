//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Sets all functions with empty status to StatusTodo
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// setInitialStatus sets all functions with empty status to StatusTodo.
func setInitialStatus(functions []model.Function) {
	for i := range functions {
		if functions[i].Status == "" {
			functions[i].Status = model.StatusTodo
		}
	}
}
