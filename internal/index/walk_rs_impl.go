//ff:func feature=index type=helper control=sequence lang=rust
//ff:what walkRsImpl: handles one impl_item — reads its implementing-type head name (rsImplTypeName, unwrapping a generic_type), pushes an rsScope{receiver} (carrying the inherited cfgTest flag so methods of a #[cfg(test)] impl are excluded), then recurses into the `body` (declaration_list) via collectRsMembers. A nameless impl still descends with an empty receiver. Mirrors walkCSTypeDecl/walkJavaTypeDecl for the impl-method case; produces the same "mod::Recv::name" qualification appendRsFunc yields.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkRsImpl pushes the impl receiver onto the scope stack and recurses into its
// body to collect method functions. cfgTest is inherited from the caller so an
// impl under #[cfg(test)] excludes its methods (cfgTestActive sees the scope).
func walkRsImpl(node *treesitter.Node, relDir string, scopes []rsScope, relPath string, out *[]model.Function, cfgTest bool) {
	recv := rsImplTypeName(node.ChildByField("type"))
	inner := make([]rsScope, 0, len(scopes)+1)
	inner = append(inner, scopes...)
	inner = append(inner, rsScope{receiver: recv, cfgTest: cfgTest})

	body := node.ChildByField("body")
	if body == nil {
		return
	}
	collectRsMembers(body, relDir, inner, relPath, out)
}
