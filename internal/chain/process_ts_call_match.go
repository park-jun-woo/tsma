//ff:func feature=chain type=implementation control=sequence
//ff:what Processes a single regex match from a TS/JS function body
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// processTSCallMatch processes a single regex match from a TS/JS function body.
func processTSCallMatch(m []string, projectRoot, file string, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	receiver := m[1]
	method := m[2]
	plainFunc := m[3]

	displayName, funcName, hasReceiver := parseTSCallNames(receiver, method, plainFunc)
	if funcName == "" {
		return
	}

	if hasReceiver && isRepoBoundary(receiver) {
		addTSBoundaryCall(displayName, "repository-interface", visited, entries)
		return
	}

	def := findTSFuncDef(projectRoot, funcName, file)
	if def == nil {
		if hasReceiver {
			addTSBoundaryCall(displayName, classifyTSBoundary(receiver), visited, entries)
		}
		return
	}

	addTSInternalCall(def, displayName, projectRoot, visited, entries, depth)
}
