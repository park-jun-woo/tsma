//ff:func feature=smell type=helper control=sequence level=review lang=java
//ff:what detectJavaReflect: TS-REFL-JV-001. Node-based — fires on any method_invocation whose `name` field is getDeclaredMethod / getDeclaredField (also the plural getDeclaredMethods/getDeclaredFields), the java.lang.reflect private-introspection escape hatch. A matching string passed as an *argument* (e.g. obj.publicMethod("getDeclaredMethod")) is a string_literal, not the invocation name node, so it never fires (false-positive zero).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// javaReflectMethods are the java.lang.reflect introspection calls that reach
// private members by name.
var javaReflectMethods = map[string]bool{
	"getDeclaredMethod":  true,
	"getDeclaredField":   true,
	"getDeclaredMethods": true,
	"getDeclaredFields":  true,
}

// detectJavaReflect finds reflective getDeclaredMethod/Field calls in a test.
func detectJavaReflect(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "method_invocation" {
			return true
		}
		name := n.ChildByField("name")
		if name == nil || !javaReflectMethods[name.Text] {
			return true
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-JV-001",
			File: path,
			Line: n.StartLine(),
			Note: name.Text + "()",
		})
		return true
	})
	return findings
}
