//ff:func feature=index type=factory control=sequence lang=python
//ff:what newPyAstIndexer: constructs the Python-configured PyAstIndexer — resolves the interpreter via resolvePython and wires the existing line-based PyIndexer as the graceful fallback. The Python analogue of newTSTreeSitterIndexer: it is the one place Python D1 config meets the indexing pipeline. Unlike TypeScript, the precise path is the built-in `ast` module (parent §7-1 lock), not tree-sitter.
package index

// newPyAstIndexer returns a PyAstIndexer wired for Python. The interpreter is
// resolved here; if it is absent the indexer's Index transparently delegates to
// the line-based PyIndexer, so behavior is identical to pre-005b in environments
// without Python.
func newPyAstIndexer() *PyAstIndexer {
	return &PyAstIndexer{
		python:   resolvePython(),
		fallback: &PyIndexer{},
	}
}
