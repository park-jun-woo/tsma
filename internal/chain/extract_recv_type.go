//ff:func feature=chain type=helper control=selection
//ff:what Extracts the receiver type name from an AST expression
package chain

import "go/ast"

// extractRecvType extracts the receiver type name from an AST expression.
func extractRecvType(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return extractRecvType(e.X)
	default:
		return "Unknown"
	}
}
