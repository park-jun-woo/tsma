//ff:func feature=index type=factory control=sequence lang=csharp
//ff:what newCSharpTreeSitterIndexer: constructs the C#-configured TreeSitterIndexer — resolves the CLI + tree-sitter-c-sharp grammar via internal/treesitter, injects the C# extractor (extractCSharpFunctions) and the isCsSource/skipCsDir filters, and wires the existing line-based CsIndexer/indexCsFile as the graceful fallback. Sibling of newJavaTreeSitterIndexer (005c template); grammar/CLI absence transparently degrades to the regex path, so behavior is identical to pre-005d without tree-sitter.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// newCSharpTreeSitterIndexer returns a TreeSitterIndexer wired for C#. The
// command + grammar are resolved here; if they are absent the indexer's Index
// transparently delegates to the line-based CsIndexer, so behavior matches
// pre-005d in environments without tree-sitter (zero regression).
func newCSharpTreeSitterIndexer() *TreeSitterIndexer {
	return &TreeSitterIndexer{
		lang:         "csharp",
		command:      treesitter.ResolveCommand(),
		grammarDir:   treesitter.ResolveGrammar("csharp"),
		fallback:     &CsIndexer{},
		fileFallback: func(relPath, absPath string) []model.Function { return indexCsFile(relPath, absPath) },
		isSource:     isCsSource,
		skipDir:      skipCsDir,
		extract:      extractCSharpFunctions,
	}
}
