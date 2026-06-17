//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what rsMacroCallNames: extracts called function names from a macro token_tree, where tree-sitter does NOT parse the body into call_expressions (the macro-non-expansion AST ceiling, parent §4). A call `foo(args)` inside `assert_eq!(...)` appears as an `identifier` leaf immediately followed by a `(...)` token_tree sibling, so this returns each such identifier's text — recovering the `add(2,3)` / `nested::double(4)` calls real Rust tests overwhelmingly wrap in assert macros, which a call_expression-only walk would miss entirely.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsMacroCallNames returns the names of functions called inside a macro
// token_tree (an identifier directly followed by a (...) token_tree).
func rsMacroCallNames(tt *treesitter.Node) []string {
	var out []string
	for i := 0; i+1 < len(tt.Children); i++ {
		if tt.Children[i].Type == "identifier" && tt.Children[i+1].Type == "token_tree" {
			out = append(out, tt.Children[i].Text)
		}
	}
	return out
}
