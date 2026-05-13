//ff:func feature=endpoint type=implementation control=sequence
//ff:what Handles Echo's e.Add(method, path, handler) pattern
package endpoint

import (
	"go/ast"
	"go/token"
	"strings"
)

func matchEchoAddRoute(call *ast.CallExpr, fset *token.FileSet, relPath string, funcDecls map[string]*ast.FuncDecl) *routeRegistration {
	if len(call.Args) < 3 {
		return nil
	}
	methodLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || methodLit.Kind != token.STRING {
		return nil
	}
	httpMethod := strings.Trim(methodLit.Value, `"`)

	pathLit, ok := call.Args[1].(*ast.BasicLit)
	if !ok || pathLit.Kind != token.STRING {
		return nil
	}
	routePath := strings.Trim(pathLit.Value, `"`)

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
