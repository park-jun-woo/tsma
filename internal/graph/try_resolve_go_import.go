//ff:func feature=graph type=helper control=sequence
//ff:what Attempts import-based resolution for a Go selector call expression
package graph

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// tryResolveGoImport attempts to resolve a selector call via import path.
func tryResolveGoImport(sel *ast.SelectorExpr, callerIdx int, methodName string, imports map[string]string, functions []model.Function, idx *funcIndex) bool {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	importPath, found := imports[ident.Name]
	if !found {
		return false
	}

	return resolveGoImportCall(importPath, methodName, callerIdx, functions, idx)
}
