//ff:func feature=chain type=implementation control=sequence
//ff:what Recursively traces Go function calls within a handler body using AST inspection
package chain

import (
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/tsma/internal/model"
)

// traceFunc recursively traces function calls within the handler body.
func traceFunc(fn *funcInfo, funcs map[string]*funcInfo, fset *token.FileSet, projectRoot string, visited map[string]bool, entries *[]model.ChainEntry) {
	if fn == nil || fn.body == nil {
		return
	}

	ast.Inspect(fn.body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		processGoCallExpr(call, funcs, fset, projectRoot, visited, entries)
		return true
	})
}
