//ff:func feature=match type=helper control=sequence lang=go
//ff:what Walks a func body collecting called identifiers paired with their receiver type
package match

import "go/ast"

// collectCalledRefs walks a function body and returns the set of calledRef
// pairs (identifier name + statically-resolved receiver type) for every
// CallExpr callee. It is the receiver-aware counterpart of collectCalledIdents:
// the name is resolved the same way (calleeName), and the receiver is resolved
// by calleeReceiver against a per-body local-variable type map (built once via
// localVarTypes). A free-function call or a method call whose receiver cannot be
// determined yields a calledRef with Receiver "" (unknown). The same name may
// appear with several receivers (e.g. T{}.M() and U{}.M()), and each distinct
// (name, receiver) pair is a separate set entry. varTypes is scoped to this
// body only: callers pass the test function's own map for the test body and a
// freshly built helper map for a helper body (helper 1-hop never reuses the
// test's variable scope).
func collectCalledRefs(body *ast.BlockStmt) map[calledRef]struct{} {
	refs := make(map[calledRef]struct{})
	if body == nil {
		return refs
	}
	varTypes := localVarTypes(body)
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		name := calleeName(call.Fun)
		if name == "" {
			return true
		}
		refs[calledRef{Name: name, Receiver: calleeReceiver(call.Fun, varTypes)}] = struct{}{}
		return true
	})
	return refs
}
