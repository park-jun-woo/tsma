//ff:func feature=chain type=implementation control=iteration dimension=1
//ff:what Extracts function declarations from a parsed Go file into the index
package chain

import (
	"go/ast"
	"go/token"
)

// collectDeclaredFuncs extracts function declarations from a parsed Go file into the index.
func collectDeclaredFuncs(f *ast.File, fset *token.FileSet, relPath, pkgDir string, funcs map[string]*funcInfo) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}

		key := buildFuncKey(fd, pkgDir)
		funcs[key] = &funcInfo{
			name:      fd.Name.Name,
			file:      relPath,
			startLine: fset.Position(fd.Pos()).Line,
			endLine:   fset.Position(fd.End()).Line,
			body:      fd.Body,
		}
	}
}
