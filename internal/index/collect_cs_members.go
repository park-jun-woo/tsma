//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what collectCSMembers: iterates the direct children of a container node (compilation_unit root, namespace body, or a type body declaration_list) and hands each to dispatchCSMember. The collectJavaMembers analogue — kept to a single loop so the per-child switch (and its nested guards) live one level down in dispatchCSMember, within the depth budget.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// collectCSMembers walks node's direct children, dispatching each to
// dispatchCSMember (which emits methods/constructors/properties and recurses
// into nested namespace/type declarations).
func collectCSMembers(node *treesitter.Node, fileNs string, scopes []csScope, relPath string, out *[]model.Function) {
	for _, c := range node.Children {
		dispatchCSMember(c, fileNs, scopes, relPath, out)
	}
}
