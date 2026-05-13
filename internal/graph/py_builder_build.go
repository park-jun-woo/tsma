//ff:func feature=graph type=implementation control=sequence
//ff:what Walks Python source files and builds the call graph via regex analysis
package graph

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Build analyzes Python source to populate callers/callees/entry_point/dead fields.
func (p *PyBuilder) Build(projectRoot string, functions []model.Function) ([]model.Function, model.GraphSummary, error) {
	result := copyFunctions(functions)
	idx := newFuncIndex(result)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipPyBuildDir(path)
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "test_") {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		analyzePyFile(path, relPath, result, idx)
		return nil
	})
	if err != nil {
		return nil, model.GraphSummary{}, err
	}

	markEntryAndDead(result, false)
	summary := buildSummary(result)
	return result, summary, nil
}
