//ff:func feature=chain type=implementation control=sequence
//ff:what Adds an external Go selector call to the chain entries
package chain

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// addExternalSelectorCall adds an external selector call to the chain entries.
func addExternalSelectorCall(call *ast.CallExpr, calleeName string, visited map[string]bool, entries *[]model.ChainEntry) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}
	displayName := ident.Name + "." + calleeName
	boundary := classifyBoundary(ident.Name)
	key := "external:" + displayName
	if visited[key] {
		return
	}
	visited[key] = true
	*entries = append(*entries, model.ChainEntry{
		Func:     displayName,
		Boundary: boundary,
	})
}
