//ff:func feature=endpoint type=implementation control=sequence
//ff:what Parses a single Go file for Chi route registrations via AST inspection
package endpoint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

func parseFileForChiRoutes(filePath, projectRoot string) ([]routeRegistration, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(projectRoot, filePath)

	funcDecls := collectFuncDecls(f)

	var regs []routeRegistration

	ast.Inspect(f, func(n ast.Node) bool {
		reg := matchChiRoute(n, fset, relPath, funcDecls)
		if reg != nil {
			regs = append(regs, *reg)
		}
		return true
	})

	return regs, nil
}
