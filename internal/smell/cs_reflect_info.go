//ff:func feature=smell type=helper control=sequence level=review lang=csharp
//ff:what detectCsReflectInfo: TS-REFL-CS-002. Node-based — fires on any variable_declaration whose `type` field is the identifier MethodInfo / PropertyInfo / FieldInfo (the typed reflection handle a test declares to dynamically Invoke/GetValue/SetValue a private member). Covers locals (`MethodInfo m = ...`) and fields (a field_declaration wraps a variable_declaration) alike. A string or comment mentioning "MethodInfo" is a string_literal / comment, never a `type` field node, so it never fires (false-positive zero). AF015 original .NET pattern.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csReflectInfoTypes are the System.Reflection handle types whose declaration
// signals a dynamic (reflective) member invocation.
var csReflectInfoTypes = map[string]bool{
	"MethodInfo":   true,
	"PropertyInfo": true,
	"FieldInfo":    true,
}

// detectCsReflectInfo finds declarations of MethodInfo/PropertyInfo/FieldInfo
// reflection handles in a test.
func detectCsReflectInfo(root *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		if n.Type != "variable_declaration" {
			return true
		}
		typeNode := n.ChildByField("type")
		if typeNode == nil || !csReflectInfoTypes[typeNode.Text] {
			return true
		}
		findings = append(findings, Finding{
			Rule: "TS-REFL-CS-002",
			File: path,
			Line: typeNode.StartLine(),
			Note: typeNode.Text,
		})
		return true
	})
	return findings
}
