//ff:type feature=index type=model
//ff:what sourceFile pairs a project-root-relative path with its absolute path — the unit collectSourceFiles emits for the tree-sitter batch indexer.
package index

// sourceFile pairs a project-root-relative path with its absolute path.
type sourceFile struct {
	rel string
	abs string
}
