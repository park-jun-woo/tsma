//ff:func feature=smell type=helper control=sequence level=review lang=csharp
//ff:what detectCsReflect: TS-REFL-CS-001. Node-based — fires on any invocation_expression whose function is a member_access_expression with `name` field GetMethod / GetField / GetProperty (the System.Reflection private-introspection escape hatch; plurals GetMethods/GetFields/GetProperties included). A matching string passed as an *argument* (e.g. StringUtils.IsBlank("GetMethod")) is a string_literal, not the invocation name node, so it never fires (false-positive zero). AF015 original .NET pattern.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csReflectMethods are the System.Reflection introspection calls that reach
// members by name.
var csReflectMethods = map[string]bool{
	"GetMethod":     true,
	"GetField":      true,
	"GetProperty":   true,
	"GetMethods":    true,
	"GetFields":     true,
	"GetProperties": true,
}

// detectCsReflect finds reflective GetMethod/GetField/GetProperty calls in a test.
func detectCsReflect(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "invocation_expression" {
			return true
		}
		fn := n.ChildByField("function")
		if fn == nil || fn.Type != "member_access_expression" {
			return true
		}
		name := fn.ChildByField("name")
		if name == nil || !csReflectMethods[name.Text] {
			return true
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-CS-001",
			File: path,
			Line: n.StartLine(),
			Note: name.Text + "()",
		})
		return true
	})
	return findings
}
