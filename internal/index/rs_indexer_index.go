//ff:func feature=index type=implementation control=sequence
//ff:what Walks the project tree and collects all Rust function declarations
package index

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Index walks the project tree and collects all Rust function declarations.
func (r *RsIndexer) Index(projectRoot string) ([]model.Function, error) {
	ignorePatterns := ParseTsmIgnore(filepath.Join(projectRoot, ".tsmignore"))
	var functions []model.Function

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(projectRoot, path)
		if MatchTsmIgnore(relPath, info.Name(), info.IsDir(), ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return skipRsDir(path)
		}
		if !isRsSource(path) {
			return nil
		}

		fileFuncs := indexRsFile(relPath, path)
		functions = append(functions, fileFuncs...)
		return nil
	})

	return functions, err
}
