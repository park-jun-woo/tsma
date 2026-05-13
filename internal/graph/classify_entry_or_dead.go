//ff:func feature=graph type=helper control=sequence
//ff:what Classifies a caller-less function as entry point or dead code
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// classifyEntryOrDead classifies a function with no callers.
func classifyEntryOrDead(fn *model.Function, goMode bool) {
	name := fn.Name

	if name == "main" || name == "init" {
		fn.EntryPoint = true
		return
	}

	if goMode && isGoTestFunc(name) {
		fn.EntryPoint = true
		return
	}

	if fn.Exported {
		fn.EntryPoint = true
		return
	}

	fn.Dead = true
}
