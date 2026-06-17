//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what walkCSTypeDecl: handles one scope-introducing declaration (block namespace or class/struct/interface/record/enum) — reads its dotted `name` (csDottedName, so `namespace A.B` becomes one A.B scope segment), pushes a csScope so members are qualified Namespace.Outer.Inner.Member, then recurses into the `body` (declaration_list / namespace body) via collectCSMembers. A nameless declaration is skipped. Mirrors walkJavaTypeDecl without descending into method bodies.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkCSTypeDecl pushes the namespace/type name onto the scope stack and recurses
// into its body to collect member functions.
func walkCSTypeDecl(node *treesitter.Node, fileNs string, scopes []csScope, relPath string, out *[]model.Function) {
	name := csDottedName(node.ChildByField("name"))
	if name == "" {
		return
	}
	inner := make([]csScope, 0, len(scopes)+1)
	inner = append(inner, scopes...)
	inner = append(inner, csScope{typeName: name})

	body := node.ChildByField("body")
	if body == nil {
		return
	}
	collectCSMembers(body, fileNs, inner, relPath, out)
}
