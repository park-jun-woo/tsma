//ff:func feature=graph type=implementation control=sequence
//ff:what Walks Go source files and builds the call graph via AST analysis
package graph

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Build analyzes Go source to populate callers/callees/entry_point/dead fields.
func (g *GoBuilder) Build(projectRoot string, functions []model.Function) ([]model.Function, model.GraphSummary, error) {
	result := copyFunctions(functions)
	idx := newFuncIndex(result)
	fset := token.NewFileSet()

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipGoBuildDir(path)
		}
		if !isGoBuildSource(path) {
			return nil
		}

		f, parseErr := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if parseErr != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		pkgDir := filepath.Dir(relPath)
		if pkgDir == "." {
			pkgDir = ""
		}

		imports := collectGoImports(f)
		analyzeGoFile(f, relPath, pkgDir, imports, result, idx)
		return nil
	})
	if err != nil {
		return nil, model.GraphSummary{}, err
	}

	markEntryAndDead(result, true)
	summary := buildSummary(result)
	return result, summary, nil
}

// isGoBuildSource checks if a path is a non-test non-mock Go source.
func isGoBuildSource(path string) bool {
	if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
		return false
	}
	return !strings.HasPrefix(filepath.Base(path), "mock_")
}
