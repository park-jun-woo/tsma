//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what csFuncFromNode: builds a model.Function from a method / constructor / destructor / property declaration node — name from the `name` field child, precise StartLine/EndLine from the node span (the multi-line-signature accuracy the line-based path lacks), QualifiedName via buildCsQualifiedName(fileNs, scopes, name) identical to appendCsFunc, Exported from csNodeExported (a `public` modifier). Returns false for a nameless declaration (e.g. an operator declaration), so it is simply not indexed.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// csFuncFromNode converts a method/constructor/property declaration node into a
// model.Function. The node span gives the exact StartLine..EndLine range feeding
// D3 branch coverage.
func csFuncFromNode(node *treesitter.Node, fileNs string, scopes []csScope, relPath string) (model.Function, bool) {
	nameNode := node.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	return model.Function{
		QualifiedName: buildCsQualifiedName(fileNs, scopes, name),
		Name:          name,
		File:          relPath,
		StartLine:     node.StartLine(),
		EndLine:       node.EndLine(),
		Exported:      csNodeExported(node),
		Status:        model.StatusTodo,
	}, true
}
