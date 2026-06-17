//ff:func feature=index type=helper control=iteration dimension=1 lang=typescript
//ff:what collectTSMethods: from a class_declaration node, appends a model.Function for each method_definition in the class body (via tsMethodFromDefinition), qualified as <pkgDir>.<Class>.<method> exactly like the line-based tryMatchTSMethod.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// collectTSMethods appends a model.Function for each method_definition in the
// class body.
func collectTSMethods(classNode *treesitter.Node, relDir, relPath string, out *[]model.Function) {
	className := ""
	if n := classNode.ChildByField("name"); n != nil {
		className = n.Text
	}
	body := classNode.ChildByType("class_body")
	if body == nil {
		return
	}
	for _, m := range body.Children {
		if fn, ok := tsMethodFromDefinition(m, className, relDir, relPath); ok {
			*out = append(*out, fn)
		}
	}
}
