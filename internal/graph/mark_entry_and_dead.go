//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Marks entry points and dead functions based on caller count
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// markEntryAndDead sets EntryPoint and Dead flags based on callers.
// goMode enables Go-specific rules (main, init, Test*, exported).
func markEntryAndDead(functions []model.Function, goMode bool) {
	for i := range functions {
		if len(functions[i].Callers) > 0 {
			continue
		}
		classifyEntryOrDead(&functions[i], goMode)
	}
}
