//ff:func feature=smell type=helper control=sequence level=review lang=typescript
//ff:what detectTSAsAny: TS-REFL-TS-001. Fires only when a member access lands on an `as any` cast — a member_expression whose object (through optional parens) is an as_expression to the `any` type. This is the runtime private-bypass cheese (TS `private` is compile-time only). A legitimate `as T` cast and a bare `as any` with no member access never fire (false-positive zero, per the plan's narrowing).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectTSAsAny finds `(x as any).member` private-bypass casts. It walks every
// member_expression and fires when the receiver is an as_expression to `any`.
func detectTSAsAny(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "member_expression" {
			return true
		}
		ae := unwrapAsExpr(n.ChildByField("object"))
		if ae == nil || !asExprIsAny(ae) {
			return true
		}
		prop := ""
		if p := n.ChildByField("property"); p != nil {
			prop = p.Text
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-TS-001",
			File: path,
			Line: n.StartLine(),
			Note: "(as any)." + prop,
		})
		return true
	})
	return findings
}
