//ff:func feature=smell type=helper control=sequence level=review lang=rust
//ff:what detectRsUnsafe: TS-REFL-RS-001. Node-based — fires on any unsafe_block (a test forcing raw memory access) or unsafe-modified function_item within a test scope. A "unsafe" word in a string literal or comment is not an unsafe_block/function_modifiers node, so it never fires (false-positive zero). Rust's escape hatch is unsafe, not reflection; scanning is scoped to test code (rsTestScopeNodes) so legitimate production/FFI unsafe is never flagged.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// detectRsUnsafe finds unsafe blocks and unsafe fns inside a test scope.
func detectRsUnsafe(scope *treesitter.Node, path string) []Finding {
	var findings []Finding
	treesitter.Walk(scope, func(n *treesitter.Node) bool {
		if n.Type == "unsafe_block" {
			findings = append(findings, Finding{Rule: "TS-REFL-RS-001", File: path, Line: n.StartLine(), Note: "unsafe block"})
		}
		if n.Type == "function_item" && rsFnIsUnsafe(n) {
			findings = append(findings, Finding{Rule: "TS-REFL-RS-001", File: path, Line: n.StartLine(), Note: "unsafe fn"})
		}
		return true
	})
	return findings
}
