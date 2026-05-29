//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Collects a test's direct callees plus same-file 1-hop helper callees
package match

import "go/ast"

// expandTestRefs returns the set of source identifier names referenced by a
// single test function: its directly called identifiers, plus the directly
// called identifiers of any same-file non-Test helper it invokes (1-hop only;
// deeper indirection is intentionally not followed).
func expandTestRefs(test *ast.FuncDecl, funcs map[string]*ast.FuncDecl) map[string]struct{} {
	refs := make(map[string]struct{})
	direct := collectCalledIdents(test.Body)
	for name := range direct {
		refs[name] = struct{}{}
		helper, ok := funcs[name]
		if !ok || isTestFuncName(name) {
			continue
		}
		mergeHelperRefs(refs, helper)
	}
	return refs
}
