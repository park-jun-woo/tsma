//ff:func feature=index type=helper control=sequence dimension=1 lang=python
//ff:what indexPyFileast: the precise per-file Python indexing step — runs the ast dump (runPyAst) and converts it (parsePyAst) into model.Functions. The ast analogue of indexPyFile; on any subprocess/parse error it returns the error so PyAstIndexer.Index falls back to the line-based indexPyFile for just that file (zero regression).
package index

import "github.com/park-jun-woo/tsma/internal/model"

// indexPyFileast extracts a single Python file's functions via the ast
// subprocess, returning an error when the interpreter or parse fails.
func indexPyFileast(relPath, absPath, python string) ([]model.Function, error) {
	data, err := runPyAst(python, pyAstDefScript, absPath)
	if err != nil {
		return nil, err
	}
	return parsePyAst(data, relPath)
}
