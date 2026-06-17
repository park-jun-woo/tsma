//ff:func feature=index type=helper control=iteration dimension=1 lang=rust
//ff:what collectRsMembers: iterates the direct children of a container node (source_file root, mod body, or impl body declaration_list), threading the #[cfg(test)] pending flag through dispatchRsMember so an attribute_item set on one child guards the next. The collectCSMembers analogue — kept to a single loop so the per-child switch (and its nested cfg guard) live one level down in dispatchRsMember, within the depth budget.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// collectRsMembers walks node's direct children, dispatching each to
// dispatchRsMember and carrying the pending #[cfg(test)] guard between siblings.
func collectRsMembers(node *treesitter.Node, relDir string, scopes []rsScope, relPath string, out *[]model.Function) {
	pending := false
	for _, c := range node.Children {
		pending = dispatchRsMember(c, relDir, scopes, relPath, out, pending)
	}
}
