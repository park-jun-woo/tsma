//ff:func feature=index type=implementation control=sequence
//ff:what Walks the project tree and collects all Python function declarations
package index

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Index walks the project tree and collects all Python function declarations.
func (p *PyIndexer) Index(projectRoot string) ([]model.Function, error) {
	var functions []model.Function

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipPyDir(path)
		}
		if !isPySource(path) {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		fileFuncs := indexPyFile(relPath, path)
		functions = append(functions, fileFuncs...)
		return nil
	})

	return functions, err
}
