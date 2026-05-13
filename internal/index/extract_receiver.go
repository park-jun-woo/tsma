//ff:func feature=index type=helper control=selection
//ff:what Extracts the receiver type name from an AST expression stripping pointer indirection
package index

import "go/ast"

// extractReceiver extracts the receiver type name from an AST expression.
func extractReceiver(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return extractReceiver(e.X)
	case *ast.IndexExpr:
		return extractReceiver(e.X)
	case *ast.IndexListExpr:
		return extractReceiver(e.X)
	default:
		return "Unknown"
	}
}
