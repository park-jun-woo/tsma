//ff:func feature=graph type=implementation control=selection
//ff:what Dispatches Go call expression resolution through 3-stage strategy
package graph

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// resolveGoCall resolves a call expression via same-package, import, or method matching.
func resolveGoCall(call *ast.CallExpr, callerIdx int, pkgDir string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		resolveGoSamePackage(fun.Name, callerIdx, pkgDir, functions, idx)
	case *ast.SelectorExpr:
		resolveGoSelectorCall(fun, callerIdx, imports, functions, idx)
	}
}
