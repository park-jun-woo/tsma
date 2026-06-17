//ff:func feature=index type=helper control=iteration dimension=1 lang=typescript
//ff:what collectTSArrowFuncs: from a lexical/variable declaration node, appends a model.Function for each declarator that tsArrowFuncFromDeclarator recognizes as bound to an arrow/function expression. Non-function consts are skipped (same filter as the line-based path).
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// collectTSArrowFuncs appends a model.Function for each function-valued
// declarator in a lexical/variable declaration.
func collectTSArrowFuncs(node *treesitter.Node, relDir, relPath string, exported bool, out *[]model.Function) {
	for _, d := range node.Children {
		if fn, ok := tsArrowFuncFromDeclarator(d, relDir, relPath, exported); ok {
			*out = append(*out, fn)
		}
	}
}
