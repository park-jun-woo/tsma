//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what csNodeExported: reports whether a method/constructor/property declaration is public by inspecting its direct `modifier` children — the node-based analogue of the line-based `strings.HasPrefix(trimmed, "public")`. Each access modifier is its own `(modifier)` leaf whose Text is the keyword ("public"/"private"/"static"/...), so a "public" substring inside an attribute, type, or identifier never produces a false positive.
package index

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csNodeExported returns true when the declaration carries a public modifier.
func csNodeExported(node *treesitter.Node) bool {
	for _, c := range node.Children {
		if c.Type == "modifier" && c.Text == "public" {
			return true
		}
	}
	return false
}
