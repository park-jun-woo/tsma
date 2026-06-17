//ff:func feature=index type=factory control=selection
//ff:what Returns the appropriate Indexer implementation for a given language string
package index

// NewIndexer returns the appropriate indexer for the given language.
func NewIndexer(lang string) Indexer {
	switch lang {
	case "go":
		return &GoIndexer{}
	case "typescript":
		return newTSTreeSitterIndexer()
	case "python":
		return newPyAstIndexer()
	case "rust":
		return &RsIndexer{}
	case "java":
		return newJavaTreeSitterIndexer()
	case "csharp":
		return newCSharpTreeSitterIndexer()
	default:
		return &UnsupportedIndexer{Lang: lang}
	}
}
