//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Extracts and resolves all call expressions from a single Python line
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolvePyCallsInLine extracts and resolves all call expressions from a line.
func resolvePyCallsInLine(line string, callerIdx int, relPath string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	matches := pyCallPattern.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		identifier := m[1]
		method := m[2]

		if pyBuiltins[identifier] && method == "" {
			continue
		}

		if method != "" {
			resolvePyMethodCall(method, callerIdx, functions, idx)
		} else {
			resolvePyBareCall(identifier, callerIdx, relPath, functions, idx)
		}
	}
}
