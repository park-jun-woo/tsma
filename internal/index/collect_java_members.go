//ff:func feature=index type=helper control=iteration dimension=1 lang=java
//ff:what collectJavaMembers: iterates the direct children of a container node (program root or a type body) and hands each to dispatchJavaMember. The walkTSTopLevel analogue — kept to a single loop so the per-child switch (and its nested guards) live one level down in dispatchJavaMember, within the depth budget.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// collectJavaMembers walks node's direct children, dispatching each to
// dispatchJavaMember (which emits methods/constructors and recurses into nested
// type declarations).
func collectJavaMembers(node *treesitter.Node, pkg string, scopes []javaScope, relPath string, out *[]model.Function) {
	for _, c := range node.Children {
		dispatchJavaMember(c, pkg, scopes, relPath, out)
	}
}
