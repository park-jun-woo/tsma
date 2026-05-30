//ff:func feature=match type=helper control=selection lang=go
//ff:what Extracts the bare receiver type name from a method declaration's receiver expr
package match

import "go/ast"

// srcReceiver extracts the bare receiver type name from a source method
// declaration's receiver type expression, stripping pointer indirection and
// generic type parameters so *T, T, and T[X] all map to "T". It mirrors
// index.extractReceiver's rule (pointer/value normalized to the same name) but
// lives in the match package — index.extractReceiver is unexported and returns
// "Unknown" on an unrecognized shape, whereas here an unrecognized receiver
// yields "" so it cannot accidentally collide with a real type name. Used by
// BuildPkgSourceReceivers to collect the set of receivers a method name is
// declared on within a package.
func srcReceiver(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return srcReceiver(e.X)
	case *ast.IndexExpr:
		return srcReceiver(e.X)
	case *ast.IndexListExpr:
		return srcReceiver(e.X)
	default:
		return ""
	}
}
