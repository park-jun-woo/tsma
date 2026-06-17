//ff:func feature=index type=helper control=iteration dimension=1
//ff:what (TreeSitterIndexer).fallbackFiles: re-indexes every collected file with the line-based fileFallback. Used when the batch parse fails outright so indexing still produces results (graceful fallback).
package index

import "github.com/park-jun-woo/tsma/internal/model"

// fallbackFiles re-indexes every collected file with the line-based path.
func (t *TreeSitterIndexer) fallbackFiles(files []sourceFile) []model.Function {
	var functions []model.Function
	for _, f := range files {
		functions = append(functions, t.fileFallback(f.rel, f.abs)...)
	}
	return functions
}
