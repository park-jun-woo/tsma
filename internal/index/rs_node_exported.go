//ff:func feature=index type=helper control=iteration dimension=1 lang=rust
//ff:what rsNodeExported: reports whether a function_item is `pub` by inspecting its direct `visibility_modifier` child — the node-based analogue of the line-based `strings.HasPrefix(trimmed, "pub")`. The visibility keyword is its own `(visibility_modifier)` leaf whose Text is "pub" (or "pub(crate)"); a "pub" substring inside an identifier, attribute, or doc comment never produces a false positive. Mirrors csNodeExported/javaNodeExported.
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsNodeExported returns true when the declaration carries a pub visibility.
func rsNodeExported(node *treesitter.Node) bool {
	for _, c := range node.Children {
		if c.Type == "visibility_modifier" && strings.HasPrefix(strings.TrimSpace(c.Text), "pub") {
			return true
		}
	}
	return false
}
