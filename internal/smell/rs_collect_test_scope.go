//ff:func feature=smell type=helper control=selection lang=rust
//ff:what rsCollectTestScope: classifies one direct child of the source_file root for test-scope collection and returns the carry-forward pending-test-attribute flag. A test-ish attribute_item (rsAttrMentionsTest) sets pending; a guarded mod_item contributes its body (the in-file #[cfg(test)] module) and a guarded function_item contributes itself (a top-level #[test] fn in an integration test); any other node preserves pending. The single place test scopes are recognized — keeps rsTestScopeNodes a flat loop.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsCollectTestScope appends test-scope nodes (a #[cfg(test)] mod body or a
// #[test] function) to out and returns the updated pending-attribute flag.
func rsCollectTestScope(c *treesitter.Node, pending bool, out *[]*treesitter.Node) bool {
	switch c.Type {
	case "attribute_item":
		return pending || rsAttrMentionsTest(c)
	case "mod_item":
		if b := c.ChildByField("body"); pending && b != nil {
			*out = append(*out, b)
		}
		return false
	case "function_item":
		if pending {
			*out = append(*out, c)
		}
		return false
	default:
		return pending
	}
}
