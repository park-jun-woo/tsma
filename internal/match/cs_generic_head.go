//ff:func feature=match type=helper control=iteration dimension=1 lang=csharp
//ff:what csGenericHead: returns the head name node of a generic_name (`Foo<Bar>` → the `Foo` identifier child), i.e. the first direct identifier / qualified_name child, skipping the type_argument_list. Returns nil when none is present. Extracted from csSimpleTypeOrName so its per-node switch stays within the depth-2 nesting budget (the generic case loop lives here).
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csGenericHead returns the head type-name node of a generic_name, or nil.
func csGenericHead(n *treesitter.Node) *treesitter.Node {
	for _, c := range n.Children {
		if c.Type == "identifier" || c.Type == "qualified_name" {
			return c
		}
	}
	return nil
}
