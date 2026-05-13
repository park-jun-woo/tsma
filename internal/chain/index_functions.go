//ff:func feature=chain type=implementation control=sequence
//ff:what Walks the project tree and indexes all Go function declarations via AST parsing
package chain

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// indexFunctions parses all Go files and builds a map of function declarations.
func indexFunctions(projectRoot string) (map[string]*funcInfo, *token.FileSet, error) {
	fset := token.NewFileSet()
	funcs := make(map[string]*funcInfo)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipGoDir(path)
		}
		if !isGoSourceFile(path) {
			return nil
		}

		f, parseErr := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if parseErr != nil {
			return nil // skip unparseable files
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		pkgDir := filepath.Dir(relPath)
		collectDeclaredFuncs(f, fset, relPath, pkgDir, funcs)
		return nil
	})

	return funcs, fset, err
}
