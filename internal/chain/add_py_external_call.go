//ff:func feature=chain type=implementation control=sequence
//ff:what Adds an external or unresolvable Python call to the chain
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// addPyExternalCall adds an external or unresolvable Python call to the chain.
func addPyExternalCall(callExpr string, imports map[string]string, visited map[string]bool, entries *[]model.ChainEntry) {
	boundary := classifyPyBoundary(callExpr, imports)
	key := "external:" + callExpr
	if visited[key] {
		return
	}
	visited[key] = true
	*entries = append(*entries, model.ChainEntry{
		Func:     callExpr,
		Boundary: boundary,
	})
}
