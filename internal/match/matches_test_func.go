//ff:func feature=match type=helper control=sequence
//ff:what Checks whether an AST declaration is a Test* function whose name contains the target
package match

import (
	"go/ast"
	"strings"
)

// matchesTestFunc checks whether an AST declaration is a Test* function
// whose name contains the target function name.
func matchesTestFunc(decl ast.Decl, funcName string) bool {
	fd, ok := decl.(*ast.FuncDecl)
	if !ok {
		return false
	}
	testName := fd.Name.Name
	if !strings.HasPrefix(testName, "Test") {
		return false
	}
	if fd.Type.Params == nil || len(fd.Type.Params.List) == 0 {
		return false
	}
	suffix := testName[len("Test"):]
	if len(funcName) == 0 {
		return false
	}
	capitalized := strings.ToUpper(funcName[:1]) + funcName[1:]
	return strings.Contains(suffix, capitalized) || strings.Contains(suffix, funcName)
}
