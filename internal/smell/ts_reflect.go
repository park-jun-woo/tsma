//ff:func feature=smell type=helper control=sequence level=review lang=typescript
//ff:what detectTSReflect: TS-REFL-TS-002. Node-based — fires on any member_expression whose object is the identifier `Reflect` (Reflect.get/set/ownKeys/...), the dynamic-introspection escape hatch. A `Reflect` substring in a string/comment is not an identifier node so it never fires (false-positive zero).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectTSReflect finds `Reflect.<x>` member accesses in a test file.
func detectTSReflect(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "member_expression" {
			return true
		}
		obj := n.ChildByField("object")
		if obj == nil || obj.Type != "identifier" || obj.Text != "Reflect" {
			return true
		}
		prop := ""
		if p := n.ChildByField("property"); p != nil {
			prop = p.Text
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-TS-002",
			File: path,
			Line: n.StartLine(),
			Note: "Reflect." + prop,
		})
		return true
	})
	return findings
}
