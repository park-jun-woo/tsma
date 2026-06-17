//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what csDottedName: resolves the dotted name from a declaration's `name` node — the leaf Text directly for a single-segment identifier (`namespace App` / `class Foo`), or the "."-joined identifier leaves of a qualified_name (`namespace A.B` → "A.B"). Returns "" for a nil node. Shared by csFileNamespace (file-scoped namespace) and walkCSTypeDecl (block namespaces and types), so both produce the same Namespace.Type.Member segments buildCsQualifiedName expects.
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// csDottedName returns the dotted name of a name node: its Text when it is a
// leaf identifier, otherwise its identifier leaves joined with ".".
func csDottedName(name *treesitter.Node) string {
	if name == nil {
		return ""
	}
	if name.Text != "" {
		return name.Text
	}
	var parts []string
	treesitter.Walk(name, func(n *treesitter.Node) bool {
		if n.Type == "identifier" && n.Text != "" {
			parts = append(parts, n.Text)
		}
		return true
	})
	return strings.Join(parts, ".")
}
