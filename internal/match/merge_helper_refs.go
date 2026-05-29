//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Merges a helper function's directly called identifiers into a ref set
package match

import "go/ast"

// mergeHelperRefs adds the directly called identifier names of a helper
// function into refs. Only the helper's own direct callees are merged (the
// 1-hop boundary); the helper's callees are not recursively expanded.
func mergeHelperRefs(refs map[string]struct{}, helper *ast.FuncDecl) {
	for name := range collectCalledIdents(helper.Body) {
		refs[name] = struct{}{}
	}
}
