//ff:func feature=smell type=helper control=sequence level=review lang=rust
//ff:what detectRsTransmute: TS-REFL-RS-002. Node-based — fires on any call_expression whose function resolves (rsCallName) to `transmute` (std::mem::transmute or a `use`-imported bare transmute), the type-system escape hatch a test uses to reinterpret bytes. There is no legitimate non-reinterpreting transmute, so the bare name is a precise target; a string literal "transmute" is a string_literal, not the call function node, so it never fires (false-positive zero).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectRsTransmute finds std::mem::transmute calls inside a test scope.
func detectRsTransmute(scope *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(scope, func(n *treesitter.Node) bool {
		if n.Type != "call_expression" {
			return true
		}
		if rsCallName(n.ChildByField("function")) != "transmute" {
			return true
		}
		findings = append(findings, Finding{Rule: "TS-REFL-RS-002", File: path, Line: n.StartLine(), Note: "transmute()"})
		return true
	})
	return findings
}
