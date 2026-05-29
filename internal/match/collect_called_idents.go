//ff:func feature=match type=helper control=sequence lang=go
//ff:what Walks a func body and collects the names of all directly called identifiers
package match

import "go/ast"

// collectCalledIdents walks a function body and returns the set of identifier
// names that appear as the callee of a CallExpr. For a plain call Foo(...) the
// name is the Ident "Foo"; for a method/selector call x.Foo(...) the name is
// the selector "Foo". Names are returned as a set (map[string]struct{}).
func collectCalledIdents(body *ast.BlockStmt) map[string]struct{} {
	names := make(map[string]struct{})
	if body == nil {
		return names
	}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if name := calleeName(call.Fun); name != "" {
			names[name] = struct{}{}
		}
		return true
	})
	return names
}
