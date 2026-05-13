//ff:func feature=chain type=helper control=selection
//ff:what Extracts the function name and selector flag from a Go call expression
package chain

import "go/ast"

// extractCalleeName extracts the function name and selector flag from a call expression.
func extractCalleeName(call *ast.CallExpr) (string, bool) {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name, false
	case *ast.SelectorExpr:
		return fun.Sel.Name, true
	default:
		return "", false
	}
}
