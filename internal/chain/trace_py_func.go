//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Recursively traces Python function calls within a function body
package chain

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// tracePyFunc recursively traces function calls within a Python function body.
func tracePyFunc(fn *pyFuncInfo, funcs map[string]*pyFuncInfo, projectRoot string, imports map[string]string, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	if fn == nil || depth >= maxPyTraceDepth {
		return
	}

	body := strings.Join(fn.bodyLines, "\n")
	calls := pyCallRe.FindAllStringSubmatch(body, -1)

	for _, call := range calls {
		processPyCallMatch(call, funcs, projectRoot, imports, visited, entries, depth)
	}
}
