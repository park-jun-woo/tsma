//ff:func feature=match type=helper control=selection lang=typescript
//ff:what tsCalleeName: resolves the bare called name from a tree-sitter callee node — the identifier for `foo(...)`, the property for `obj.foo(...)`. The TS analogue of Go's calleeName; member calls collapse to the trailing name so a source method is found regardless of receiver expression.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// tsCalleeName returns the bare callee name of a call/new target node, or "".
func tsCalleeName(fn *treesitter.Node) string {
	switch fn.Type {
	case "identifier":
		return fn.Text
	case "member_expression":
		if p := fn.ChildByField("property"); p != nil {
			return p.Text
		}
	}
	return ""
}
