//ff:func feature=smell type=helper control=selection lang=rust
//ff:what rsCallName: resolves the rightmost called name from a call_expression's `function` node — the leaf Text for a bare identifier (`transmute(x)`), the `name` field for a scoped_identifier (`std::mem::transmute` → "transmute"), or the `field` field for a field_expression (`v.as_ptr()` → handled by the method detector). Returns "" for an unexpected shape. The node-based name extraction the transmute detector keys on, so a string literal "transmute" never matches.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsCallName returns the simple invoked name from a call_expression function
// node, or "".
func rsCallName(fn *treesitter.Node) string {
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
	default:
		return ""
	}
}
