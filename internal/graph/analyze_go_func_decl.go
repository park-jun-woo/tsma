//ff:func feature=graph type=implementation control=sequence
//ff:what Inspects a single Go function body to extract and resolve call expressions
package graph

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeGoFuncDecl inspects a single function body for call expressions.
func analyzeGoFuncDecl(fd *ast.FuncDecl, pkgDir string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	callerQN := buildGoQualifiedName(fd, pkgDir)
	callerIdx, found := idx.byQualified[callerQN]
	if !found {
		return
	}

	ast.Inspect(fd.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		resolveGoCall(call, callerIdx, pkgDir, imports, functions, idx)
		return true
	})
}
