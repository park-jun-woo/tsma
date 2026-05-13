//ff:func feature=graph type=helper control=sequence
//ff:what Resolves a Go selector call via import path lookup then method name matching
package graph

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// resolveGoSelectorCall resolves a selector expression (pkg.Func or obj.Method).
func resolveGoSelectorCall(sel *ast.SelectorExpr, callerIdx int, imports map[string]string, functions []model.Function, idx *funcIndex) {
	methodName := sel.Sel.Name

	if tryResolveGoImport(sel, callerIdx, methodName, imports, functions, idx) {
		return
	}

	resolveGoMethodCall(methodName, callerIdx, functions, idx)
}
