//ff:func feature=chain type=helper control=sequence
//ff:what Adds an external/boundary TS call entry if not already visited
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// addTSBoundaryCall adds an external/boundary call entry if not already visited.
func addTSBoundaryCall(displayName, boundary string, visited map[string]bool, entries *[]model.ChainEntry) {
	key := "external:" + displayName
	if visited[key] {
		return
	}
	visited[key] = true
	*entries = append(*entries, model.ChainEntry{
		Func:     displayName,
		Boundary: boundary,
	})
}
