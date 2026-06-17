//ff:func feature=index type=helper control=iteration dimension=1 lang=typescript
//ff:what walkTSTopLevel: iterates a node's direct children at program/export scope, handing each to walkTSChild. It never descends into a function body, matching the Go indexer's top-level-only semantics.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkTSTopLevel dispatches every direct child of node (program or
// export_statement) to walkTSChild.
func walkTSTopLevel(node *treesitter.Node, relDir, relPath string, exported bool, out *[]model.Function) {
	for _, c := range node.Children {
		walkTSChild(c, relDir, relPath, exported, out)
	}
}
