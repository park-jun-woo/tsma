//ff:func feature=match type=helper control=selection lang=csharp
//ff:what csInvokedName: resolves the bare method name from an invocation_expression's `function` node — the `name` field for a member_access_expression (`obj.Foo()` → "Foo", also `obj.Foo<T>()`), the leaf Text for a bare identifier (`Foo()`), or the first identifier child for a generic_name (`Foo<T>()`). Returns "" when none is found. The C# analogue of the method-name half of collectJavaCalledNames, keeping `obj.Foo()` attributable to a method named Foo regardless of receiver or type arguments.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csInvokedName returns the simple invoked method name from a function node, or "".
func csInvokedName(fn *treesitter.Node) string {
	if fn == nil {
		return ""
	}
	switch fn.Type {
	case "member_access_expression":
		return csSimpleTypeOrName(fn.ChildByField("name"))
	default:
		return csSimpleTypeOrName(fn)
	}
}
