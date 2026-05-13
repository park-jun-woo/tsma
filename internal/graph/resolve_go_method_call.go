//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Resolves a Go method call by name with ambiguity handling for multiple candidates
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolveGoMethodCall resolves a method call by name with ambiguity handling.
func resolveGoMethodCall(methodName string, callerIdx int, functions []model.Function, idx *funcIndex) {
	candidates := idx.byName[methodName]
	var matched []int
	for _, ci := range candidates {
		if ci != callerIdx {
			matched = append(matched, ci)
		}
	}

	if len(matched) == 0 {
		return
	}

	ambiguous := len(matched) > 1
	for _, ci := range matched {
		addEdge(functions, callerIdx, ci, ambiguous)
	}
}
