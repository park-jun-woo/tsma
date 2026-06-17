//ff:func feature=match type=helper control=selection lang=rust
//ff:what rsInvokedName: resolves the bare invoked name from a call_expression's `function` node — the leaf Text for a bare identifier (`free_fn()`), the `name` field for a scoped_identifier (`mod::free_fn()` or `Calculator::new()` → the trailing segment), or the `field` field for a field_expression (`obj.method()` → "method"). Returns "" for an unexpected shape. The Rust analogue of csInvokedName, keeping a function attributable by its own name regardless of receiver/path.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsInvokedName returns the simple invoked name from a call_expression function
// node, or "".
func rsInvokedName(fn *treesitter.Node) string {
	if fn == nil {
		return ""
	}
	switch fn.Type {
	case "identifier":
		return fn.Text
	case "scoped_identifier":
		if name := fn.ChildByField("name"); name != nil {
			return name.Text
		}
		return ""
	case "field_expression":
		if field := fn.ChildByField("field"); field != nil {
			return field.Text
		}
		return ""
	default:
		return ""
	}
}
