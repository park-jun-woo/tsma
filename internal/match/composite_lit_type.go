//ff:func feature=match type=helper control=selection lang=go
//ff:what Resolves a receiver expression to a bare type name via composite literals only
package match

import "go/ast"

// compositeLitType resolves a call's receiver expression x (the .X of a
// SelectorExpr call) to a bare type name, but only for the two supported
// syntactic patterns: a composite literal used directly, optionally addressed.
//
//	T{...}.M()        -> "T"   (CompositeLit, Type is Ident)
//	(&T{...}).M()     -> "T"   (UnaryExpr & of CompositeLit)
//
// Surrounding parentheses are unwrapped (the (&T{...}) form parses as a
// ParenExpr). Pointer/value is normalized to the same bare type name (the
// leading & is dropped), mirroring index.extractReceiver's pointer-stripping
// rule. Any other
// receiver expression — a variable, a constructor call, a field access, an
// interface argument, a selector type (pkg.T) on the literal, etc. — yields ""
// (unknown), which the matching policy treats conservatively. Local-variable
// bindings are resolved separately by localVarTypes; this helper only inspects
// the expression in front of the dot.
func compositeLitType(x ast.Expr) string {
	switch e := x.(type) {
	case *ast.ParenExpr:
		return compositeLitType(e.X)
	case *ast.UnaryExpr:
		if e.Op.String() == "&" {
			return compositeLitType(e.X)
		}
		return ""
	case *ast.CompositeLit:
		if id, ok := e.Type.(*ast.Ident); ok {
			return id.Name
		}
		return ""
	default:
		return ""
	}
}
