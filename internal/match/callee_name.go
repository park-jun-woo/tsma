//ff:func feature=match type=helper control=selection lang=go
//ff:what Extracts the bare identifier name from a call expression's callee
package match

import "go/ast"

// calleeName returns the bare identifier name of a call expression's callee.
// Foo(...)   -> "Foo"   (Ident)
// x.Foo(...) -> "Foo"   (SelectorExpr method name)
// Anything else (e.g. call returned by another call) yields "".
func calleeName(fun ast.Expr) string {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	default:
		return ""
	}
}
