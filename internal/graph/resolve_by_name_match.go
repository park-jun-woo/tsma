//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Matches a function name against all candidates and adds edges with ambiguity
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolveByNameMatch matches a function by name and adds edges with ambiguity handling.
func resolveByNameMatch(funcName string, callerIdx int, functions []model.Function, idx *funcIndex) {
	candidates := idx.byName[funcName]
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
