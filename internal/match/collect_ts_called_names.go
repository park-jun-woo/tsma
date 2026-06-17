//ff:func feature=match type=helper control=sequence lang=typescript
//ff:what collectTSCalledNames: walks a parsed test file and returns the set of names invoked as call/new targets — the TS analogue of collectCalledIdents. Test-harness noise (expect/describe/it) is collected too but harmlessly, since lookup is keyed by the source function's own name.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// collectTSCalledNames returns the set of bare names called (call_expression) or
// constructed (new_expression) anywhere in the tree.
func collectTSCalledNames(root *treesitter.Node) map[string]struct{} {
	names := make(map[string]struct{})
	treesitter.Walk(root, func(n *treesitter.Node) bool {
		switch n.Type {
		case "call_expression":
			if fn := n.ChildByField("function"); fn != nil {
				if name := tsCalleeName(fn); name != "" {
					names[name] = struct{}{}
				}
			}
		case "new_expression":
			if c := n.ChildByField("constructor"); c != nil {
				if name := tsCalleeName(c); name != "" {
					names[name] = struct{}{}
				}
			}
		}
		return true
	})
	return names
}
