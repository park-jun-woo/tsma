//ff:type feature=index type=model
//ff:what xmlTreeBuilder: the streaming state for ParseXML — the running node stack, the current source's root/name, and the accumulated Sources. Holding it in one struct keeps the per-token handlers flat (depth ≤2) and lets ParseXML stay a simple token loop.
package treesitter

// xmlTreeBuilder accumulates Sources from a stream of tree-sitter XML tokens.
type xmlTreeBuilder struct {
	sources       []Source
	stack         []*Node
	curRoot       *Node
	curSourceName string
	inSources     bool
}
