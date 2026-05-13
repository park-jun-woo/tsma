//ff:func feature=chain type=implementation control=sequence
//ff:what Traces the Python call chain from a handler by indexing functions and recursing
package chain

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Trace traces the call chain starting from the handler function.
func (t *PyTracer) Trace(projectRoot string, handler model.FuncLocation) ([]model.ChainEntry, error) {
	// Build an index of all Python functions in the project.
	funcs, err := indexPyFunctions(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("index python functions: %w", err)
	}

	// Find the handler function.
	handlerFunc := findPyFuncAtLocation(funcs, handler.File, handler.StartLine)
	if handlerFunc == nil {
		return nil, nil
	}

	// Collect imports from the handler's file to classify external calls.
	handlerImports := collectImports(filepath.Join(projectRoot, handler.File))

	// Trace recursively.
	visited := make(map[string]bool)
	var entries []model.ChainEntry
	tracePyFunc(handlerFunc, funcs, projectRoot, handlerImports, visited, &entries, 0)

	return entries, nil
}
