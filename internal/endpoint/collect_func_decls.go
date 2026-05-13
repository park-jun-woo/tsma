//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Collects all function declarations from a Go AST file into a map
package endpoint

import "go/ast"

func collectFuncDecls(f *ast.File) map[string]*ast.FuncDecl {
	funcDecls := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcDecls[fd.Name.Name] = fd
	}
	return funcDecls
}
