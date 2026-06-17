//ff:func feature=smell type=helper control=sequence level=review lang=java
//ff:what detectJavaSetAccessible: TS-REFL-JV-002. Node-based — fires on a method_invocation named setAccessible whose argument list contains a `true` literal (the private-access-force escape hatch). setAccessible(false) — restoring the guard — carries a `false` node and never fires, so the rule targets only the cheese (false-positive zero).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectJavaSetAccessible finds setAccessible(true) calls in a test.
func detectJavaSetAccessible(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "method_invocation" {
			return true
		}
		name := n.ChildByField("name")
		if name == nil || name.Text != "setAccessible" {
			return true
		}
		if !javaArgsHaveTrue(n.ChildByField("arguments")) {
			return true
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-JV-002",
			File: path,
			Line: n.StartLine(),
			Note: "setAccessible(true)",
		})
		return true
	})
	return findings
}
