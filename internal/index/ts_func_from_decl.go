//ff:func feature=index type=helper control=sequence lang=typescript
//ff:what tsFuncFromDecl: builds a model.Function from a tree-sitter function_declaration node — name from the `name` field child, precise StartLine/EndLine from the node span (the multi-line-signature accuracy the line-based path lacks). Exported is carried from the enclosing export_statement.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// tsFuncFromDecl converts a function_declaration node into a model.Function. It
// returns false when the declaration has no usable name. The node span gives the
// exact body range (StartLine..EndLine) feeding D3 branch coverage.
func tsFuncFromDecl(node *treesitter.Node, relDir, relPath string, exported bool) (model.Function, bool) {
	nameNode := node.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	return model.Function{
		QualifiedName: buildQualifiedName(relDir, "", name),
		Name:          name,
		File:          relPath,
		StartLine:     node.StartLine(),
		EndLine:       node.EndLine(),
		Exported:      exported,
		Status:        model.StatusTodo,
	}, true
}
