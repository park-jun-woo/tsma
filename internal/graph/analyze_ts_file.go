//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Reads a TS/JS file and resolves call edges for functions defined in it
package graph

import (
	"bufio"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeTSFile reads a TS/JS file and resolves call edges for functions defined in it.
func analyzeTSFile(absPath, relPath string, functions []model.Function, idx *funcIndex) {
	f, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	imports := collectTSImports(lines)

	for callerIdx := range functions {
		if functions[callerIdx].File != relPath {
			continue
		}
		analyzeTSFuncBody(lines, callerIdx, relPath, imports, functions, idx)
	}
}
