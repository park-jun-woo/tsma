//ff:func feature=match type=helper control=sequence lang=go
//ff:what Maps local variables bound once to a composite-literal type within a func body
package match

import "go/ast"

// localVarTypes scans a function body and returns a map from local variable
// name to the bare type name it was bound to, but only for the simplest
// statically-resolvable bindings:
//
//	f := &T{...}   -> f: "T"
//	f := T{...}    -> f: "T"
//
// Only single-LHS-Ident assignments whose RHS is a composite literal
// (optionally addressed) are tracked. A variable that is assigned (`:=` or `=`)
// more than once — re-binding, reassignment, or assignment inside a branch or
// loop — is dropped (mapped to "" and excluded), so a later f.M() call resolves
// to unknown rather than to a possibly-stale type. This is the conservative
// "prefer non-detection over mis-attribution" rule: any ambiguity yields no
// receiver. The walk is purely syntactic (go/ast); no type checking.
func localVarTypes(body *ast.BlockStmt) map[string]string {
	types := make(map[string]string)
	poisoned := make(map[string]struct{})
	if body == nil {
		return types
	}
	ast.Inspect(body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		// Only single LHS = single RHS bindings are tracked; anything else
		// (multi-assign, tuple unpack) poisons every LHS name it touches.
		if len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			for _, lhs := range as.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					poison(types, poisoned, id.Name)
				}
			}
			return true
		}
		id, ok := as.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}
		if id.Name == "_" {
			return true
		}
		// A second assignment to the same name (re-binding or reassignment)
		// makes its type ambiguous -> drop it.
		if _, seen := types[id.Name]; seen {
			poison(types, poisoned, id.Name)
			return true
		}
		if _, dead := poisoned[id.Name]; dead {
			return true
		}
		t := compositeLitType(as.Rhs[0])
		if t == "" {
			poison(types, poisoned, id.Name)
			return true
		}
		types[id.Name] = t
		return true
	})
	return types
}
