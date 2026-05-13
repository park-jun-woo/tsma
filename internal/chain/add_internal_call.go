//ff:func feature=chain type=implementation control=sequence
//ff:what Adds an internal Go function call to the chain and recurses into it
package chain

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/tsma/internal/model"
)

// addInternalCall adds an internal function call to the chain and recurses.
func addInternalCall(call *ast.CallExpr, callee *funcInfo, calleeName string, isSelector bool, funcs map[string]*funcInfo, fset *token.FileSet, projectRoot string, visited map[string]bool, entries *[]model.ChainEntry) {
	key := callee.file + ":" + fmt.Sprintf("%d", callee.startLine)
	if visited[key] {
		return
	}
	visited[key] = true

	displayName := resolveDisplayName(call, calleeName, isSelector)
	*entries = append(*entries, model.ChainEntry{
		Func:      displayName,
		File:      callee.file,
		StartLine: callee.startLine,
		EndLine:   callee.endLine,
	})

	traceFunc(callee, funcs, fset, projectRoot, visited, entries)
}
