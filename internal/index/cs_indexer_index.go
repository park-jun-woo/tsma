//ff:func feature=index type=implementation control=sequence lang=csharp
//ff:what Walks the project tree and collects all C# method declarations
package index

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Index walks the project tree and collects all C# method declarations.
func (c *CsIndexer) Index(projectRoot string) ([]model.Function, error) {
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
			return skipCsDir(path)
		}
		if !isCsSource(path) {
			return nil
		}

		fileFuncs := indexCsFile(relPath, path)
		functions = append(functions, fileFuncs...)
		return nil
	})

	return functions, err
}
