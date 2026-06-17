//ff:func feature=index type=factory control=sequence lang=rust
//ff:what newRsTreeSitterIndexer: constructs the Rust-configured TreeSitterIndexer — resolves the CLI + tree-sitter-rust grammar via internal/treesitter, injects the Rust extractor (extractRustFunctions) and the isRsSource/skipRsDir filters, and wires the existing line-based RsIndexer/indexRsFile as the graceful fallback. Sibling of newCSharpTreeSitterIndexer (005d template); grammar/CLI absence transparently degrades to the regex path, so behavior is identical to pre-005e without tree-sitter (zero regression).
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// newRsTreeSitterIndexer returns a TreeSitterIndexer wired for Rust. The command
// + grammar are resolved here; if they are absent the indexer's Index
// transparently delegates to the line-based RsIndexer, so behavior matches
// pre-005e in environments without tree-sitter (zero regression).
func newRsTreeSitterIndexer() *TreeSitterIndexer {
	return &TreeSitterIndexer{
		lang:         "rust",
		command:      treesitter.ResolveCommand(),
		grammarDir:   treesitter.ResolveGrammar("rust"),
		fallback:     &RsIndexer{},
		fileFallback: func(relPath, absPath string) []model.Function { return indexRsFile(relPath, absPath) },
		isSource:     isRsSource,
		skipDir:      skipRsDir,
		extract:      extractRustFunctions,
	}
}
