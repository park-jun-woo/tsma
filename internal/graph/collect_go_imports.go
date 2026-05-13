//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Extracts import alias to import path mapping from a Go file AST
package graph

import (
	"go/ast"
	"strings"
)

// collectGoImports extracts import alias -> import path mapping from a Go file.
func collectGoImports(f *ast.File) map[string]string {
	imports := make(map[string]string)
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			alias = parts[len(parts)-1]
		}
		if alias != "_" && alias != "." {
			imports[alias] = path
		}
	}
	return imports
}
