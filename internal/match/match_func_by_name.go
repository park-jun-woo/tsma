//ff:func feature=match type=helper control=sequence lang=go
//ff:what Looks up a function's test references by its model.Function.Name key
package match

import "github.com/park-jun-woo/tsma/internal/model"

// MatchFuncByName returns the test references attributed to a function's bare
// name within its package, using a prebuilt index. It is the lookup entry point
// keyed by model.Function.Name (consistent with the indexer's bare names).
func MatchFuncByName(idx *PkgTestIndex, fn *model.Function) ([]testRef, bool) {
	if idx == nil || fn == nil {
		return nil, false
	}
	return idx.Lookup(fn.Name)
}
