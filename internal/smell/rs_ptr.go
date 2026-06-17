//ff:func feature=smell type=helper control=sequence level=review lang=rust
//ff:what detectRsPtr: TS-REFL-RS-003. Node-based — fires on a std::ptr::read/write scoped call (rsPtrScopedCall) or an as_ptr/as_mut_ptr field access (rsPtrFieldAccess) within a test scope, the raw-pointer escape hatch a test uses to read/write memory bypassing the type system. Both predicates require the precise `ptr`/`as_*_ptr` shape, so a plain `reader.read()` or a string literal never matches (false-positive zero). The per-node classification lives in the two predicate helpers to keep this a flat walk.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectRsPtr finds std::ptr read/write calls and as_ptr extractions in a test
// scope, the TS-REFL-RS-003 raw-pointer escape hatch.
func detectRsPtr(scope *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(scope, func(n *treesitter.Node) bool {
		switch n.Type {
		case "call_expression":
			if note, ok := rsPtrScopedCall(n.ChildByField("function")); ok {
				findings = append(findings, Finding{Rule: "TS-REFL-RS-003", File: path, Line: n.StartLine(), Note: note})
			}
		case "field_expression":
			if note, ok := rsPtrFieldAccess(n.ChildByField("field")); ok {
				findings = append(findings, Finding{Rule: "TS-REFL-RS-003", File: path, Line: n.StartLine(), Note: note})
			}
		}
		return true
	})
	return findings
}
