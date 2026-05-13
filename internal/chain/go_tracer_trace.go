//ff:func feature=chain type=implementation control=sequence
//ff:what Traces the Go call chain from a handler by indexing functions and recursing
package chain

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Trace traces the call chain starting from the handler function.
func (t *GoTracer) Trace(projectRoot string, handler model.FuncLocation) ([]model.ChainEntry, error) {
	// Parse all Go files in the project to build a function index.
	funcs, fset, err := indexFunctions(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("index functions: %w", err)
	}

	// Find the handler function.
	handlerFunc := findFuncAtLocation(funcs, handler.File, handler.StartLine)
	if handlerFunc == nil {
		return nil, nil
	}

	// Trace the call chain recursively.
	visited := make(map[string]bool)
	var entries []model.ChainEntry
	traceFunc(handlerFunc, funcs, fset, projectRoot, visited, &entries)

	return entries, nil
}
