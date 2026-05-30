//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Merges a helper function's directly called identifiers into a ref set
package match

import "go/ast"

// mergeHelperRefs adds the directly called identifier references (name +
// receiver) of a helper function into refs. Only the helper's own direct
// callees are merged (the 1-hop boundary); the helper's callees are not
// recursively expanded. Receivers are resolved against the helper's own local
// variable map (collectCalledRefs builds it from the helper body), never the
// caller test's scope, so a composite literal constructed inside the helper is
// attributed to its type while a value the helper receives as an argument
// resolves to unknown.
func mergeHelperRefs(refs map[calledRef]struct{}, helper *ast.FuncDecl) {
	for ref := range collectCalledRefs(helper.Body) {
		refs[ref] = struct{}{}
	}
}
