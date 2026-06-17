//ff:func feature=index type=helper control=sequence lang=rust
//ff:what rsFuncFromNode: builds a model.Function from a function_item node — name from the `name` field child, precise StartLine/EndLine from the node span (the multi-line-signature accuracy the line-based path lacks), QualifiedName via buildRsQualifiedName(relDir, scopes, name) identical to appendRsFunc, Exported from rsNodeExported (a `pub` visibility_modifier). Returns false for a nameless declaration. Mirrors csFuncFromNode/javaFuncFromNode.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsFuncFromNode converts a function_item node into a model.Function. The node
// span gives the exact StartLine..EndLine range feeding D3 branch coverage.
func rsFuncFromNode(node *treesitter.Node, relDir string, scopes []rsScope, relPath string) (model.Function, bool) {
	nameNode := node.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	return model.Function{
		QualifiedName: buildRsQualifiedName(relDir, scopes, name),
		Name:          name,
		File:          relPath,
		StartLine:     node.StartLine(),
		EndLine:       node.EndLine(),
		Exported:      rsNodeExported(node),
		Status:        model.StatusTodo,
	}, true
}
