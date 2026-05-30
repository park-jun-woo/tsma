//ff:func feature=match type=helper control=sequence lang=go
//ff:what Resolves a method call's receiver type from composite literals or local vars
package match

import "go/ast"

// calleeReceiver returns the statically-resolved receiver type name for a call
// expression's callee, or "" (unknown) when it cannot be determined. It only
// applies to selector calls x.M(); a plain Ident call Foo() is a free function
// and always returns "". For x.M() it resolves x by two patterns, in order:
//
//  1. a composite literal in front of the dot — T{...}.M() or (&T{...}).M()
//     (via compositeLitType);
//  2. a bare local variable whose single binding was a composite literal,
//     looked up in varTypes (f := &T{...}; f.M()).
//
// Pointer/value is normalized to the same bare type name. Any other receiver
// expression (constructor return, field access, interface argument, multi-bound
// variable, selector type) yields "", and the matching policy then treats the
// call's receiver as unknown.
func calleeReceiver(fun ast.Expr, varTypes map[string]string) string {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	if t := compositeLitType(sel.X); t != "" {
		return t
	}
	if id, ok := sel.X.(*ast.Ident); ok {
		if t, ok := varTypes[id.Name]; ok {
			return t
		}
	}
	return ""
}
