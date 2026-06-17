//ff:type feature=index type=implementation
//ff:what TreeSitterIndexer: the language-neutral, grammar-injected Indexer that drives `tree-sitter parse` (via internal/treesitter) and extracts model.Functions from the parse tree. It is the precise path laid on top of the line-based indexers; when the CLI/grammar is absent it delegates to `fallback` (graceful fallback, zero regression). 005a wires it for TypeScript; 005b~e reuse it by supplying their own grammar + extractor + fallback.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// TreeSitterIndexer indexes a project by shelling out to the tree-sitter CLI and
// interpreting the parse tree with a language-specific extractor. Everything
// except `extract` (and the source/skip filters) is language-neutral, so the
// subprocess + XML pipeline is shared across languages.
type TreeSitterIndexer struct {
	// lang is the language label (for diagnostics).
	lang string
	// command is the tree-sitter executable (resolved name or absolute path).
	command string
	// grammarDir is the tree-sitter grammar directory passed via -p; when "" the
	// CLI is expected to resolve the grammar itself (configured parser dirs).
	grammarDir string
	// fallback is the whole-project line-based indexer used when tree-sitter is
	// unavailable.
	fallback Indexer
	// fileFallback re-indexes a single file with the line-based path when one
	// file fails to parse but the batch otherwise succeeds.
	fileFallback func(relPath, absPath string) []model.Function
	// isSource selects indexable source files (e.g. isTSSource).
	isSource func(path string) bool
	// skipDir prunes directories during the walk (e.g. skipTSDir).
	skipDir func(path string) error
	// extract interprets a parsed file's root node into model.Functions.
	extract func(root *treesitter.Node, relPath, relDir string) []model.Function
}
