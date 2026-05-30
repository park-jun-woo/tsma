//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Collects a test's direct callees plus same-file 1-hop helper callees
package match

import "go/ast"

// expandTestRefs returns the set of source identifier references (name +
// statically-resolved receiver type) made by a single test function: its
// directly called identifiers, plus the directly called identifiers of any
// same-file non-Test helper it invokes (1-hop only; deeper indirection is
// intentionally not followed). Receivers are resolved per body: the test body
// uses its own local-variable map and each helper uses its own (a helper never
// inherits the test's variable scope), so a value passed into a helper as an
// argument resolves to unknown receiver — conservative by design.
func expandTestRefs(test *ast.FuncDecl, funcs map[string]*ast.FuncDecl) map[calledRef]struct{} {
	refs := make(map[calledRef]struct{})
	direct := collectCalledRefs(test.Body)
	for ref := range direct {
		refs[ref] = struct{}{}
		helper, ok := funcs[ref.Name]
		if !ok || isTestFuncName(ref.Name) {
			continue
		}
		mergeHelperRefs(refs, helper)
	}
	return refs
}
