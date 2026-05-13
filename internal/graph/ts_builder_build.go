//ff:func feature=graph type=implementation control=sequence
//ff:what Walks TS/JS source files and builds the call graph via regex analysis
package graph

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Build analyzes TS/JS source to populate callers/callees/entry_point/dead fields.
func (t *TSBuilder) Build(projectRoot string, functions []model.Function) ([]model.Function, model.GraphSummary, error) {
	result := copyFunctions(functions)
	idx := newFuncIndex(result)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipTSBuildDir(path)
		}
		if !isTSBuildFile(path) {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		analyzeTSFile(path, relPath, result, idx)
		return nil
	})
	if err != nil {
		return nil, model.GraphSummary{}, err
	}

	markEntryAndDead(result, false)
	summary := buildSummary(result)
	return result, summary, nil
}

// isTSBuildFile returns true if the path is a TS/JS source file for analysis.
func isTSBuildFile(path string) bool {
	if strings.HasSuffix(path, ".d.ts") {
		return false
	}
	if strings.HasSuffix(path, ".test.ts") || strings.HasSuffix(path, ".test.js") {
		return false
	}
	if strings.HasSuffix(path, ".spec.ts") || strings.HasSuffix(path, ".spec.js") {
		return false
	}
	return strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".js")
}
