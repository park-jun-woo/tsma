//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Computes GraphSummary from the annotated function list
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// buildSummary computes graph summary statistics from annotated functions.
func buildSummary(functions []model.Function) model.GraphSummary {
	var s model.GraphSummary
	s.Nodes = len(functions)
	for i := range functions {
		s.Edges += len(functions[i].Callees)
		if functions[i].EntryPoint {
			s.EntryPoints++
		}
		if functions[i].Dead {
			s.Dead++
		}
	}
	return s
}
