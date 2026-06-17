//ff:func feature=index type=helper control=iteration dimension=1 lang=java
//ff:what javaNodeExported: reports whether a method/constructor declaration is public by inspecting its `modifiers` child — the node-based analogue of the line-based `strings.HasPrefix(trimmed, "public")`. The modifiers node carries the literal `public` keyword as its direct character text (annotations like @Override are separate child elements), so a "public" substring inside an annotation or type never produces a false positive.
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// javaNodeExported returns true when the declaration carries a public modifier.
func javaNodeExported(node *treesitter.Node) bool {
	mods := node.ChildByType("modifiers")
	if mods == nil {
		return false
	}
	for _, field := range strings.Fields(mods.Text) {
		if field == "public" {
			return true
		}
	}
	return false
}
