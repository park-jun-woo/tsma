//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Extracts and resolves all call expressions from a single TS/JS line
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// resolveTSCallsInLine extracts and resolves all call expressions from a line.
func resolveTSCallsInLine(line string, callerIdx int, relPath string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	matches := tsCallPattern.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		identifier := m[1]
		method := m[2]

		if isTSBuiltin(identifier) {
			continue
		}

		if method != "" {
			resolveTSMethodCall(method, callerIdx, functions, idx)
		} else {
			resolveTSBareCall(identifier, callerIdx, relPath, functions, idx)
		}
	}
}
