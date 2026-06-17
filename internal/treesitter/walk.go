//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Walk: depth-first pre-order traversal of a tree-sitter Node. Returning false from the visitor prunes that node's subtree. Used by node-based detectors (smell D4) that must inspect every member_expression/call regardless of nesting.
package treesitter

// Walk visits node and all descendants depth-first, pre-order. If visit returns
// false the node's children are skipped. A nil node is a no-op.
func Walk(node *Node, visit func(*Node) bool) {
	if node == nil {
		return
	}
	if !visit(node) {
		return
	}
	for _, c := range node.Children {
		Walk(c, visit)
	}
}
