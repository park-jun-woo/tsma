//ff:func feature=index type=helper control=sequence lang=rust
//ff:what emitRsFunc: appends a function_item as a model.Function unless it is test-only — the cfgTestActive(scopes, pending) guard (the in-file #[cfg(test)] mod exclusion) plus the rsFuncFromNode conversion, factored out of dispatchRsMember so that switch's function_item case stays within the depth budget (Q1 ≤2). A cfg(test)-guarded or nameless declaration is simply not emitted.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// emitRsFunc appends c to out as a model.Function unless it is inside a
// #[cfg(test)] scope (cfgTestActive) or nameless.
func emitRsFunc(c *treesitter.Node, relDir string, scopes []rsScope, relPath string, out *[]model.Function, pending bool) {
	if cfgTestActive(scopes, pending) {
		return
	}
	if fn, ok := rsFuncFromNode(c, relDir, scopes, relPath); ok {
		*out = append(*out, fn)
	}
}
