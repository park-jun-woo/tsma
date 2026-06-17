//ff:func feature=index type=helper control=iteration dimension=1
//ff:what (Node).ChildByType: returns the first direct child of the given grammar node type (e.g. "class_body", "as_expression"), or nil — used when the relationship is by type rather than field.
package treesitter

// ChildByType returns the first direct child with the given node type, or nil.
func (n *Node) ChildByType(nodeType string) *Node {
	for _, c := range n.Children {
		if c.Type == nodeType {
			return c
		}
	}
	return nil
}
