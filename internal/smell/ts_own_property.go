//ff:func feature=smell type=helper control=sequence level=review lang=typescript
//ff:what detectTSOwnProperty: TS-REFL-TS-003. Node-based — fires on Object.getOwnPropertyNames / Object.getOwnPropertyDescriptor (reaching private/non-enumerable members). Public iteration (Object.keys/values/entries) is deliberately excluded, so normal object traversal never fires (false-positive zero).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// tsOwnPropertyReach are the Object.* methods that reach past public iteration
// into own (incl. non-enumerable / private-by-convention) members.
var tsOwnPropertyReach = map[string]bool{
	"getOwnPropertyNames":      true,
	"getOwnPropertyDescriptor": true,
}

// detectTSOwnProperty finds Object.getOwnPropertyNames/Descriptor accesses.
func detectTSOwnProperty(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "member_expression" {
			return true
		}
		obj := n.ChildByField("object")
		if obj == nil || obj.Type != "identifier" || obj.Text != "Object" {
			return true
		}
		prop := n.ChildByField("property")
		if prop == nil || !tsOwnPropertyReach[prop.Text] {
			return true
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-TS-003",
			File: path,
			Line: n.StartLine(),
			Note: "Object." + prop.Text,
		})
		return true
	})
	return findings
}
