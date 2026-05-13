//ff:func feature=graph type=implementation control=iteration dimension=1
//ff:what Walks all function bodies in a Go file and resolves call edges
package graph

import (
	"go/ast"

	"github.com/park-jun-woo/tsma/internal/model"
)

// analyzeGoFile walks all function bodies in a file and resolves call edges.
func analyzeGoFile(f *ast.File, relPath, pkgDir string, imports map[string]string, functions []model.Function, idx *funcIndex) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		analyzeGoFuncDecl(fd, pkgDir, imports, functions, idx)
	}
}
