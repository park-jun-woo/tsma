//ff:type feature=index type=implementation lang=python
//ff:what PyAstIndexer: the precise Python indexer (Phase005b D1). It shells out to the built-in `ast` module (`python -c`) to extract def/async def/methods/nested functions with exact line ranges, laid on top of the line-based PyIndexer. When no Python interpreter is on PATH it delegates wholesale to `fallback` (the line-based PyIndexer), so behavior is identical to pre-005b in environments without Python. ast.parse only parses (never imports), so D1 has no import side-effect risk.
package index

// PyAstIndexer indexes Python functions via the built-in `ast` module run as a
// subprocess. python is the resolved interpreter ("" when absent), and fallback
// is the line-based PyIndexer used when the interpreter is unavailable or a
// single file fails to parse.
type PyAstIndexer struct {
	python   string
	fallback Indexer
}
