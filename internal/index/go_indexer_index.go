//ff:func feature=index type=implementation control=sequence
//ff:what Walks the project tree and collects all Go function declarations via AST
package index

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Index walks the project tree and collects all Go function declarations.
func (g *GoIndexer) Index(projectRoot string) ([]model.Function, error) {
	fset := token.NewFileSet()
	var functions []model.Function

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipGoDir(path)
		}
		if !isGoSource(path) {
			return nil
		}

		f, parseErr := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if parseErr != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		pkgDir := pkgDirOf(relPath)
		collectGoFunctions(f, fset, relPath, pkgDir, &functions)

		return nil
	})

	return functions, err
}

// collectGoFunctions extracts function declarations from a parsed Go file.
func collectGoFunctions(f *ast.File, fset *token.FileSet, relPath, pkgDir string, functions *[]model.Function) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		fn := buildGoFunction(fd, fset, relPath, pkgDir)
		*functions = append(*functions, fn)
	}
}
