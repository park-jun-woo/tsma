//ff:func feature=match type=helper control=sequence lang=rust
//ff:what collectRsCalledNames: walks a parsed Rust subtree and returns the set of called function names — both ordinary call_expressions (bare `foo()`, scoped `m::foo()`/`Type::new()`, or method `obj.foo()`, via rsInvokedName) AND calls hidden inside macro token_trees (via rsMacroCallNames), because real Rust tests wrap most calls in `assert_eq!`/`assert!` which tree-sitter leaves as raw tokens. The Rust analogue of collectCsCalledNames; lookup is keyed by the source function's own name so any collected extras are harmless.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// collectRsCalledNames returns the set of bare function names invoked anywhere
// under node — through call_expressions and through macro token_trees.
func collectRsCalledNames(node *treesitter.Node) map[string]struct{} {
	names := make(map[string]struct{})
	treesitter.Walk(node, func(n *treesitter.Node) bool {
		switch n.Type {
		case "call_expression":
			if name := rsInvokedName(n.ChildByField("function")); name != "" {
				names[name] = struct{}{}
			}
		case "token_tree":
			for _, name := range rsMacroCallNames(n) {
				names[name] = struct{}{}
			}
		}
		return true
	})
	return names
}
