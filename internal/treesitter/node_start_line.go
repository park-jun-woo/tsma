//ff:func feature=index type=helper control=sequence
//ff:what (Node).StartLine: the node's 1-based start line (0-based grammar srow + 1).
package treesitter

// StartLine returns the 1-based start line of the node.
func (n *Node) StartLine() int { return n.SRow + 1 }
