//ff:func feature=graph type=helper control=sequence
//ff:what Resolves a bare TS/JS function call by same-directory then name matching
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolveTSBareCall resolves a bare function call.
func resolveTSBareCall(funcName string, callerIdx int, callerFile string, functions []model.Function, idx *funcIndex) {
	callerDir := pkgDirOf(callerFile)
	var targetQN string
	if callerDir == "" {
		targetQN = funcName
	} else {
		targetQN = callerDir + "." + funcName
	}

	if calleeIdx, found := idx.byQualified[targetQN]; found && calleeIdx != callerIdx {
		addEdge(functions, callerIdx, calleeIdx, false)
		return
	}

	resolveByNameMatch(funcName, callerIdx, functions, idx)
}
