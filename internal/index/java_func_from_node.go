//ff:func feature=index type=helper control=sequence lang=java
//ff:what javaFuncFromNode: builds a model.Function from a method_declaration / constructor_declaration node — name from the `name` field child, precise StartLine/EndLine from the node span (the multi-line-signature accuracy the line-based path lacks), QualifiedName via buildJavaQualifiedName(pkg, scopes, name) identical to appendJavaFunc, Exported from javaNodeExported (a `public` modifier). Returns false for a nameless declaration.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// javaFuncFromNode converts a method/constructor declaration node into a
// model.Function. The node span gives the exact StartLine..EndLine range feeding
// D3 branch coverage.
func javaFuncFromNode(node *treesitter.Node, pkg string, scopes []javaScope, relPath string) (model.Function, bool) {
	nameNode := node.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	return model.Function{
		QualifiedName: buildJavaQualifiedName(pkg, scopes, name),
		Name:          name,
		File:          relPath,
		StartLine:     node.StartLine(),
		EndLine:       node.EndLine(),
		Exported:      javaNodeExported(node),
		Status:        model.StatusTodo,
	}, true
}
