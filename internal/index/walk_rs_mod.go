//ff:func feature=index type=helper control=sequence lang=rust
//ff:what walkRsMod: handles one inline mod_item — reads its `name` identifier, pushes an rsScope{module} (carrying the inherited cfgTest flag, so a #[cfg(test)] mod excludes every function inside it — the in-file unit-test module the plan requires be NOT indexed), then recurses into the `body` (declaration_list) via collectRsMembers. A nameless or body-less mod is skipped (a `mod foo;` file-module declaration has no inline body). Mirrors walkCSTypeDecl for the namespace case.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkRsMod pushes the module name onto the scope stack and recurses into its
// body. cfgTest is inherited so a #[cfg(test)] mod excludes its functions.
func walkRsMod(node *treesitter.Node, relDir string, scopes []rsScope, relPath string, out *[]model.Function, cfgTest bool) {
	nameNode := node.ChildByField("name")
	body := node.ChildByField("body")
	if nameNode == nil || body == nil {
		return
	}
	inner := make([]rsScope, 0, len(scopes)+1)
	inner = append(inner, scopes...)
	inner = append(inner, rsScope{module: nameNode.Text, cfgTest: cfgTest})
	collectRsMembers(body, relDir, inner, relPath, out)
}
