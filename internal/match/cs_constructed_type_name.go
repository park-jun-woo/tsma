//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what csConstructedTypeName: resolves the simple type name from an object_creation_expression's `type` node — the text directly for a bare identifier (`new Foo()`), or the first identifier descendant for a generic_name / qualified_name (`new Foo<Bar>()`, `new Ns.Foo()`). Returns "" when none is found. Keeps `new Foo()` attributable to a constructor named Foo regardless of type arguments. Delegates to csSimpleTypeOrName, the shared simple-name resolver also used by csInvokedName.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csConstructedTypeName returns the simple constructed type name, or "".
func csConstructedTypeName(typeNode *treesitter.Node) string {
	return csSimpleTypeOrName(typeNode)
}
