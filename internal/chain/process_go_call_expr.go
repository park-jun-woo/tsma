//ff:func feature=chain type=implementation control=sequence
//ff:what Processes a single Go call expression during chain tracing
package chain

import (
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/tsma/internal/model"
)

// processGoCallExpr processes a single Go call expression during chain tracing.
func processGoCallExpr(call *ast.CallExpr, funcs map[string]*funcInfo, fset *token.FileSet, projectRoot string, visited map[string]bool, entries *[]model.ChainEntry) {
	calleeName, isSelector := extractCalleeName(call)
	if calleeName == "" {
		return
	}

	callee := findUnvisitedCallee(funcs, calleeName, visited)
	if callee != nil {
		addInternalCall(call, callee, calleeName, isSelector, funcs, fset, projectRoot, visited, entries)
		return
	}

	if isSelector {
		addExternalSelectorCall(call, calleeName, visited, entries)
	}
}
