//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Parses Go test file AST to check if any Test* name contains the target function name
package match

import (
	"go/parser"
	"go/token"
)

// containsTestFor parses a Go test file and returns true if any Test*
// function name contains the target function name.
func containsTestFor(testFilePath string, funcName string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, testFilePath, nil, 0)
	if err != nil {
		return false
	}

	for _, decl := range f.Decls {
		if matchesTestFunc(decl, funcName) {
			return true
		}
	}

	return false
}
