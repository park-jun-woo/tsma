//ff:func feature=index type=helper control=selection lang=rust
//ff:what rsImplTypeName: extracts the implementing-type name from an impl_item's `type` field — the head-extraction helper that flattens the generic wrapper so `impl Foo`, `impl Foo<T>`, and `impl Trait for Foo<T>` all yield "Foo" (the receiver segment buildRsQualifiedName joins, identical to the line-based rsImplPattern capture). A plain `type_identifier` gives its Text; a `generic_type` gives its inner `type` field's Text. Returns "" for an unexpected shape so the impl is still descended without a bogus receiver.
package index

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsImplTypeName returns the base type name of an impl block's implementing type,
// unwrapping a generic_type (Foo<T>) to its head type_identifier (Foo).
func rsImplTypeName(typeNode *treesitter.Node) string {
	if typeNode == nil {
		return ""
	}
	switch typeNode.Type {
	case "type_identifier":
		return typeNode.Text
	case "generic_type":
		if inner := typeNode.ChildByField("type"); inner != nil {
			return inner.Text
		}
		return ""
	default:
		return typeNode.Text
	}
}
