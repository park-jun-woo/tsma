//ff:func feature=index type=helper control=sequence
//ff:what (Node).EndLine: the node's 1-based end line (0-based grammar erow + 1).
package treesitter

// EndLine returns the 1-based end line of the node.
func (n *Node) EndLine() int { return n.ERow + 1 }
