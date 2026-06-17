//ff:func feature=smell type=helper control=sequence lang=typescript
//ff:what unwrapAsExpr: returns the as_expression inside an object node, unwrapping a single parenthesized_expression, or nil when the object is not an as-cast — the receiver-unwrap step of the TS-REFL-TS-001 detector.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// unwrapAsExpr returns the as_expression in obj, unwrapping a single
// parenthesized_expression, or nil when obj is not an as-cast.
func unwrapAsExpr(obj *treesitter.Node) *treesitter.Node {
	if obj == nil {
		return nil
	}
	if obj.Type == "parenthesized_expression" {
		return obj.ChildByType("as_expression")
	}
	if obj.Type == "as_expression" {
		return obj
	}
	return nil
}
