//ff:func feature=chain type=implementation control=sequence
//ff:what Processes a single regex call match during Python chain tracing
package chain

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// processPyCallMatch processes a single regex call match during Python chain tracing.
func processPyCallMatch(call []string, funcs map[string]*pyFuncInfo, projectRoot string, imports map[string]string, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	callExpr := call[1]
	parts := strings.Split(callExpr, ".")
	funcName := parts[len(parts)-1]

	if isPyBuiltin(funcName) {
		return
	}

	callee := findPyCallee(funcs, funcName, parts)
	if callee != nil {
		addPyInternalCall(callee, callExpr, funcs, projectRoot, imports, visited, entries, depth)
		return
	}

	addPyExternalCall(callExpr, imports, visited, entries)
}
