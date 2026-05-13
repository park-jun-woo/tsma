//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Creates a deep copy of the function slice with fresh caller/callee slices
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// copyFunctions creates a deep copy of the function slice with fresh slices.
func copyFunctions(functions []model.Function) []model.Function {
	result := make([]model.Function, len(functions))
	for i := range functions {
		result[i] = functions[i]
		result[i].Callers = nil
		result[i].Callees = nil
	}
	return result
}
