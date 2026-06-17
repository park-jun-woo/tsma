//ff:func feature=index type=helper control=sequence lang=java
//ff:what walkJavaTypeDecl: handles one type declaration (class/interface/enum/record/annotation) — reads its `name` field, pushes a javaScope so members are qualified pkg.Outer.Inner.method, then recurses into the `body` (class_body/interface_body/enum_body) via collectJavaMembers. A nameless type is skipped. Mirrors the line-based brace-scope stack without descending into method bodies.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkJavaTypeDecl pushes the type's name onto the scope stack and recurses into
// its body to collect member functions.
func walkJavaTypeDecl(node *treesitter.Node, pkg string, scopes []javaScope, relPath string, out *[]model.Function) {
	nameNode := node.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return
	}
	inner := make([]javaScope, 0, len(scopes)+1)
	inner = append(inner, scopes...)
	inner = append(inner, javaScope{typeName: nameNode.Text})

	body := node.ChildByField("body")
	if body == nil {
		return
	}
	collectJavaMembers(body, pkg, inner, relPath, out)
}
