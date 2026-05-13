//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Reads a Python file and resolves call edges for functions defined in it
package graph

import (
	"bufio"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzePyFile reads a Python file and resolves call edges for functions defined in it.
func analyzePyFile(absPath, relPath string, functions []model.Function, idx *funcIndex) {
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

	imports := collectPyImports(lines)

	for callerIdx := range functions {
		if functions[callerIdx].File != relPath {
			continue
		}
		analyzePyFuncBody(lines, callerIdx, relPath, imports, functions, idx)
	}
}
