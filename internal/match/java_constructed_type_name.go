//ff:func feature=match type=helper control=sequence lang=java
//ff:what javaConstructedTypeName: resolves the simple type name from an object_creation_expression's `type` node — the text directly for a bare type_identifier (`new Foo()`), or the first type_identifier descendant for a generic_type (`new Foo<Bar>()`). Returns "" when none is found. Keeps `new Foo()` attributable to a constructor named Foo regardless of type arguments.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// javaConstructedTypeName returns the simple constructed type name, or "".
func javaConstructedTypeName(typeNode *treesitter.Node) string {
	if typeNode == nil {
		return ""
	}
	if typeNode.Type == "type_identifier" && typeNode.Text != "" {
		return typeNode.Text
	}
	found := ""
	treesitter.Walk(typeNode, func(n *treesitter.Node) bool {
		if found != "" {
			return false
		}
		if n.Type == "type_identifier" && n.Text != "" {
			found = n.Text
			return false
		}
		return true
	})
	return found
}
