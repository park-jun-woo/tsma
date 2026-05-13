//ff:func feature=chain type=implementation control=sequence
//ff:what Adds an internal Python function call to the chain and recurses
package chain

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// addPyInternalCall adds an internal Python function call to the chain and recurses.
func addPyInternalCall(callee *pyFuncInfo, callExpr string, funcs map[string]*pyFuncInfo, projectRoot string, imports map[string]string, visited map[string]bool, entries *[]model.ChainEntry, depth int) {
	key := callee.file + ":" + fmt.Sprintf("%d", callee.startLine)
	if visited[key] {
		return
	}
	visited[key] = true

	*entries = append(*entries, model.ChainEntry{
		Func:      callExpr,
		File:      callee.file,
		StartLine: callee.startLine,
		EndLine:   callee.endLine,
	})

	tracePyFunc(callee, funcs, projectRoot, imports, visited, entries, depth+1)
}
