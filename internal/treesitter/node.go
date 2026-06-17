//ff:type feature=index type=model
//ff:what Language-neutral tree-sitter parse node (one element of `tree-sitter parse --xml` output). Shared by every consumer of the tree-sitter subprocess pipeline (index D1, smell D4, and 005b~e); rows are the raw 0-based grammar rows, and the StartLine/EndLine helpers add 1 for 1-based source lines.
package treesitter

// Node is one node of a tree-sitter parse tree as emitted by
// `tree-sitter parse --xml`. It is deliberately language-neutral: the subprocess
// and XML-parsing machinery are shared, and only per-language interpreters read
// the node types. Rows are 0-based; callers add 1 for 1-based source lines.
type Node struct {
	// Type is the grammar node type, e.g. "function_declaration".
	Type string
	// Field is the field name on the parent ("name", "value", ...) or "".
	Field string
	// SRow/SCol/ERow/ECol are the 0-based start/end positions.
	SRow, SCol, ERow, ECol int
	// Text is the concatenated direct character data (set for leaf nodes such
	// as identifier/property_identifier; empty for interior nodes).
	Text string
	// Children are the direct child nodes in source order.
	Children []*Node
}
