//ff:func feature=match type=helper control=selection lang=csharp
//ff:what csSimpleTypeOrName: collapses a C# name/type node to its simple identifier — the leaf Text for a bare identifier, the trailing segment of a qualified_name (`Ns.Foo` → "Foo", via its `name` field), or the head identifier of a generic_name (`Foo<Bar>` → "Foo", its first identifier child, never a type argument). Returns "" for anything else. Shared by csInvokedName and csConstructedTypeName so invocation and construction resolve names identically.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csSimpleTypeOrName reduces a name/type node to its simple identifier, or "".
func csSimpleTypeOrName(n *treesitter.Node) string {
	if n == nil {
		return ""
	}
	switch n.Type {
	case "identifier":
		return n.Text
	case "qualified_name":
		return csSimpleTypeOrName(n.ChildByField("name"))
	case "generic_name":
		return csSimpleTypeOrName(csGenericHead(n))
	default:
		if n.Text != "" {
			return n.Text
		}
		return ""
	}
}
