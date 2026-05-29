//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Parses a Go test file and returns its top-level function declarations by name
package match

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// parseTestFileFuncs parses a Go source file and returns its top-level function
// declarations keyed by function name. Methods are keyed by their bare method
// name; on a name collision (same-package same-named methods on different
// receivers) the first declaration wins, which is acceptable because helper
// expansion only widens the referenced-identifier set.
func parseTestFileFuncs(absPath string) (map[string]*ast.FuncDecl, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, absPath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	funcs := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		if _, exists := funcs[fd.Name.Name]; !exists {
			funcs[fd.Name.Name] = fd
		}
	}
	return funcs, nil
}
