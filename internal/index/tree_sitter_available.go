//ff:func feature=index type=helper control=selection
//ff:what (TreeSitterIndexer).available: reports whether the precise tree-sitter path can run (the CLI was resolved). The gate for graceful fallback — when false, Index delegates entirely to the line-based indexer.
package index

// available reports whether the tree-sitter CLI was resolved and the precise
// path can be attempted. Grammar resolution is allowed to be deferred to the
// CLI's own parser-dir config, so only the command is required here.
func (t *TreeSitterIndexer) available() bool {
	return t.command != ""
}
