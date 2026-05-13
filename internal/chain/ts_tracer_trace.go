//ff:func feature=chain type=implementation control=sequence
//ff:what Reads handler body and delegates to traceTSFunc for recursive call extraction
package chain

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Trace traces the call chain starting from the handler function.
func (t *TSTracer) Trace(projectRoot string, handler model.FuncLocation) ([]model.ChainEntry, error) {
	if handler.File == "" {
		return nil, nil
	}

	absFile := filepath.Join(projectRoot, handler.File)
	data, err := os.ReadFile(absFile)
	if err != nil {
		return nil, nil // handler file not found; skip gracefully
	}

	lines := strings.Split(string(data), "\n")

	visited := make(map[string]bool)
	var entries []model.ChainEntry

	traceTSFunc(projectRoot, handler.File, lines, handler.StartLine, handler.EndLine, visited, &entries, 0)

	return entries, nil
}
