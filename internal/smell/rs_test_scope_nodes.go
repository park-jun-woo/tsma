//ff:func feature=smell type=helper control=iteration dimension=1 lang=rust
//ff:what rsTestScopeNodes: returns the subtree roots that constitute test code in a parsed Rust file — the bodies of in-file #[cfg(test)] modules plus any top-level #[test] functions (integration tests, whose whole file is a test crate). The escape-hatch detectors run only over these scopes so legitimate production unsafe/FFI never produces a finding (false-positive zero). Threads the pending #[test]/#[cfg(test)] attribute through rsCollectTestScope.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsTestScopeNodes returns the test-scope subtree roots (#[cfg(test)] mod bodies
// and top-level #[test] functions) of a parsed Rust source_file.
func rsTestScopeNodes(root *treesitter.Node) []*treesitter.Node {
	if root == nil {
		return nil
	}
	var scopes []*treesitter.Node
	pending := false
	for _, c := range root.Children {
		pending = rsCollectTestScope(c, pending, &scopes)
	}
	return scopes
}
