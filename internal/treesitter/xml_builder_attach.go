//ff:func feature=index type=helper control=sequence
//ff:what (xmlTreeBuilder).attach: links a new node under the current stack top, or records it as the source root when the stack is empty (the program node).
package treesitter

// attach links node under the current stack top, or records it as the source
// root when the stack is empty.
func (b *xmlTreeBuilder) attach(node *Node) {
	if len(b.stack) > 0 {
		parent := b.stack[len(b.stack)-1]
		parent.Children = append(parent.Children, node)
		return
	}
	b.curRoot = node
}
