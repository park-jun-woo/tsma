//ff:func feature=index type=implementation control=iteration dimension=1 lang=python
//ff:what (PyAstIndexer).Index: the precise Python indexing path. When no interpreter is on PATH it delegates wholesale to the line-based PyIndexer (zero regression); otherwise it walks the tree exactly like PyIndexer.Index (same .tsmignore/skipPyDir/isPySource filters) but extracts each file with the ast subprocess, per-file-falling-back to the line indexPyFile for any file ast failed to parse. Output is the same []model.Function the fallback yields, so matcher/coverage are unchanged.
package index

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Index collects Python function declarations using the built-in ast module,
// falling back to the line-based indexer when no interpreter is present or a
// file fails to parse.
func (p *PyAstIndexer) Index(projectRoot string) ([]model.Function, error) {
	if p.python == "" {
		return p.fallback.Index(projectRoot)
	}

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
			return skipPyDir(path)
		}
		if !isPySource(path) {
			return nil
		}

		fns, aerr := indexPyFileast(relPath, path, p.python)
		if aerr != nil {
			fns = indexPyFile(relPath, path)
		}
		functions = append(functions, fns...)
		return nil
	})

	return functions, err
}
