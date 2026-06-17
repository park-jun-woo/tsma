//ff:type feature=index type=model
//ff:what Source pairs one parsed file's path (as passed to the tree-sitter CLI) with its parse-tree root Node — the per-file unit ParseXML emits.
package treesitter

// Source pairs a parsed file's path (as passed to the CLI) with its root node.
type Source struct {
	Name string
	Root *Node
}
