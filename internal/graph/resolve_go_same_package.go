//ff:func feature=graph type=helper control=sequence
//ff:what Resolves a bare function call within the same Go package
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolveGoSamePackage resolves a bare function call within the same package.
func resolveGoSamePackage(funcName string, callerIdx int, pkgDir string, functions []model.Function, idx *funcIndex) {
	var targetQN string
	if pkgDir == "" {
		targetQN = funcName
	} else {
		targetQN = pkgDir + "." + funcName
	}

	if calleeIdx, found := idx.byQualified[targetQN]; found && calleeIdx != callerIdx {
		addEdge(functions, callerIdx, calleeIdx, false)
	}
}
