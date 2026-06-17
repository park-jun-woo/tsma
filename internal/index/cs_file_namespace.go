//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what csFileNamespace: extracts the dotted namespace from a file-scoped namespace declaration (`namespace A.B;`) at the top of a parsed C# file by finding the file_scoped_namespace_declaration child and joining its identifier leaf texts with "." (so a qualified_name A→B yields "A.B"). Returns "" when the file has no file-scoped namespace (block-scoped namespaces are handled as scopes by dispatchCSMember instead). Feeds buildCsQualifiedName exactly as the line-based fileNs field does.
package index

import "github.com/park-jun-woo/tsma/internal/treesitter"

// csFileNamespace returns the dotted file-scoped namespace declared at the top
// of the file, or "".
func csFileNamespace(root *treesitter.Node) string {
	ns := root.ChildByType("file_scoped_namespace_declaration")
	if ns == nil {
		return ""
	}
	return csDottedName(ns.ChildByField("name"))
}
