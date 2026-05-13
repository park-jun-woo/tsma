//ff:func feature=endpoint type=helper control=selection
//ff:what Extracts the function name from an AST expression
package endpoint

import "go/ast"

// extractFuncName extracts the function name from an expression.
func extractFuncName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.CallExpr:
		return extractFuncName(e.Fun)
	default:
		return ""
	}
}
