//ff:func feature=match type=helper control=sequence lang=java
//ff:what collectJavaCalledNames: walks a parsed JUnit test file and returns the set of names invoked as method calls (method_invocation `name` field) or constructed (object_creation_expression `type` field) anywhere in the tree — the Java analogue of collectTSCalledNames/collectCalledIdents. Test-harness noise (assertEquals/assertTrue) is collected too but harmlessly, since lookup is keyed by the source function's own name.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// collectJavaCalledNames returns the set of bare method names invoked
// (method_invocation) or type names constructed (object_creation_expression).
func collectJavaCalledNames(root *treesitter.Node) map[string]struct{} {
	names := make(map[string]struct{})
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		switch n.Type {
		case "method_invocation":
			if name := n.ChildByField("name"); name != nil && name.Text != "" {
				names[name.Text] = struct{}{}
			}
		case "object_creation_expression":
			if name := javaConstructedTypeName(n.ChildByField("type")); name != "" {
				names[name] = struct{}{}
			}
		}
		return true
	})
	return names
}
