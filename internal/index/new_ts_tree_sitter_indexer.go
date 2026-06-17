//ff:func feature=index type=factory control=sequence lang=typescript
//ff:what newTSTreeSitterIndexer: constructs the TypeScript-configured TreeSitterIndexer — resolves the CLI + grammar via internal/treesitter, injects the TS extractor and the isTSSource/skipTSDir filters, and wires the existing line-based TSIndexer/indexTSFile as the graceful fallback. This is the one place TS-specific config meets the language-neutral pipeline; 005b~e add a sibling constructor per language.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// newTSTreeSitterIndexer returns a TreeSitterIndexer wired for TypeScript/
// JavaScript. The command + grammar are resolved here; if they are absent the
// indexer's Index transparently delegates to the line-based TSIndexer, so
// behavior is identical to pre-005a in environments without tree-sitter.
func newTSTreeSitterIndexer() *TreeSitterIndexer {
	return &TreeSitterIndexer{
		lang:         "typescript",
		command:      treesitter.ResolveCommand(),
		grammarDir:   treesitter.ResolveGrammar("typescript"),
		fallback:     &TSIndexer{},
		fileFallback: func(relPath, absPath string) []model.Function { return indexTSFile(relPath, absPath) },
		isSource:     isTSSource,
		skipDir:      skipTSDir,
		extract:      extractTSFunctions,
	}
}
