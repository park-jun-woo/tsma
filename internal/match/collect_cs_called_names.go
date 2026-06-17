//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what collectCsCalledNames: walks a parsed C# test file and returns the set of names invoked as method calls (invocation_expression — member access `obj.Foo()`, bare `Foo()`, or generic `Foo<T>()`, resolved by csInvokedName) or constructed (object_creation_expression `type`, resolved by csConstructedTypeName) anywhere in the tree — the C# analogue of collectJavaCalledNames. Test-harness noise (Assert.Equal/True) is collected too but harmlessly, since lookup is keyed by the source function's own name.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// collectCsCalledNames returns the set of bare method names invoked
// (invocation_expression) or type names constructed (object_creation_expression).
func collectCsCalledNames(root *treesitter.Node) map[string]struct{} {
	names := make(map[string]struct{})
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		switch n.Type {
		case "invocation_expression":
			if name := csInvokedName(n.ChildByField("function")); name != "" {
				names[name] = struct{}{}
			}
		case "object_creation_expression":
			if name := csConstructedTypeName(n.ChildByField("type")); name != "" {
				names[name] = struct{}{}
			}
		}
		return true
	})
	return names
}
