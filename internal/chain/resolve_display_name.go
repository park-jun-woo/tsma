//ff:func feature=chain type=helper control=sequence
//ff:what Builds the display name for a Go call, including receiver if selector
package chain

import "go/ast"

// resolveDisplayName builds the display name for a call, including receiver if selector.
func resolveDisplayName(call *ast.CallExpr, calleeName string, isSelector bool) string {
	if !isSelector {
		return calleeName
	}
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name + "." + calleeName
		}
	}
	return calleeName
}
