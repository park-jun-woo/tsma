//ff:func feature=index type=helper control=iteration dimension=1
//ff:what (Node).ChildByField: returns the first direct child carrying the given tree-sitter field name (e.g. "name", "value", "object"), or nil — the primary accessor language interpreters use to pull a node's named parts.
package treesitter

// ChildByField returns the first direct child with the given field name, or nil.
func (n *Node) ChildByField(field string) *Node {
	for _, c := range n.Children {
		if c.Field == field {
			return c
		}
	}
	return nil
}
