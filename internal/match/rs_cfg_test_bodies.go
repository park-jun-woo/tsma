//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what rsCfgTestBodies: returns the declaration_list bodies of every in-file #[cfg(test)] module in a parsed Rust source file — the regions whose call sites constitute the file's own unit tests. Threads the pending #[cfg(test)] attribute through rsCollectCfgTestBody. Used by BuildRsTestIndex to attribute a source function to its file only when the in-file test module actually calls it (content-aware), not merely because a test module exists.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsCfgTestBodies returns the bodies of all #[cfg(test)] modules at the top level
// of a parsed Rust source_file.
func rsCfgTestBodies(root *treesitter.Node) []*treesitter.Node {
	if root == nil {
		return nil
	}
	var bodies []*treesitter.Node
	pending := false
	for _, c := range root.Children {
		pending = rsCollectCfgTestBody(c, pending, &bodies)
	}
	return bodies
}
