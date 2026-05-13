//ff:func feature=endpoint type=implementation control=sequence
//ff:what Checks if an AST node is a Chi route call and extracts route registration
package endpoint

import (
	"go/ast"
	"go/token"
	"strings"
)

func matchChiRoute(n ast.Node, fset *token.FileSet, relPath string, funcDecls map[string]*ast.FuncDecl) *routeRegistration {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	methodName := sel.Sel.Name

	httpMethod, isChi := chiHTTPMethods[methodName]
	if !isChi {
		return nil
	}
	if len(call.Args) < 2 {
		return nil
	}

	pathLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || pathLit.Kind != token.STRING {
		return nil
	}
	routePath := strings.Trim(pathLit.Value, `"`)

	if httpMethod == "" {
		httpMethod = methodName
	}

	lastArg := call.Args[len(call.Args)-1]
	handlerName := extractFuncName(lastArg)
	if handlerName == "" {
		return nil
	}

	startLine, endLine := 0, 0
	handlerFile := relPath
	if fd, found := funcDecls[handlerName]; found {
		startLine = fset.Position(fd.Pos()).Line
		endLine = fset.Position(fd.End()).Line
	}

	return &routeRegistration{
		method:    httpMethod,
		path:      routePath,
		handler:   handlerName,
		file:      handlerFile,
		startLine: startLine,
		endLine:   endLine,
	}
}
