//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Builds lookup indices from a function slice for fast name resolution
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// newFuncIndex builds lookup indices from a function slice.
func newFuncIndex(functions []model.Function) *funcIndex {
	idx := &funcIndex{
		byQualified: make(map[string]int, len(functions)),
		byName:      make(map[string][]int, len(functions)),
	}
	for i := range functions {
		idx.byQualified[functions[i].QualifiedName] = i
		idx.byName[functions[i].Name] = append(idx.byName[functions[i].Name], i)
	}
	return idx
}
