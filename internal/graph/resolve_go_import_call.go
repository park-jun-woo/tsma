//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Resolves a Go call through an imported package path to project functions
package graph

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// resolveGoImportCall resolves a call through an imported package.
func resolveGoImportCall(importPath, methodName string, callerIdx int, functions []model.Function, idx *funcIndex) bool {
	candidates := idx.byName[methodName]
	for _, ci := range candidates {
		if ci == callerIdx {
			continue
		}
		funcPkgDir := pkgDirOf(functions[ci].File)
		if strings.HasSuffix(importPath, funcPkgDir) || funcPkgDir == importPath {
			addEdge(functions, callerIdx, ci, false)
			return true
		}
	}
	return false
}
