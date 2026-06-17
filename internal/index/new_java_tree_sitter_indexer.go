//ff:func feature=index type=factory control=sequence lang=java
//ff:what newJavaTreeSitterIndexer: constructs the Java-configured TreeSitterIndexer — resolves the CLI + tree-sitter-java grammar via internal/treesitter, injects the Java extractor (extractJavaFunctions) and the isJavaSource/skipJavaDir filters, and wires the existing line-based JavaIndexer/indexJavaFile as the graceful fallback. Sibling of newTSTreeSitterIndexer (005a template); grammar/CLI absence transparently degrades to the regex path, so behavior is identical to pre-005c without tree-sitter.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// newJavaTreeSitterIndexer returns a TreeSitterIndexer wired for Java. The
// command + grammar are resolved here; if they are absent the indexer's Index
// transparently delegates to the line-based JavaIndexer, so behavior matches
// pre-005c in environments without tree-sitter (zero regression).
func newJavaTreeSitterIndexer() *TreeSitterIndexer {
	return &TreeSitterIndexer{
		lang:         "java",
		command:      treesitter.ResolveCommand(),
		grammarDir:   treesitter.ResolveGrammar("java"),
		fallback:     &JavaIndexer{},
		fileFallback: func(relPath, absPath string) []model.Function { return indexJavaFile(relPath, absPath) },
		isSource:     isJavaSource,
		skipDir:      skipJavaDir,
		extract:      extractJavaFunctions,
	}
}
