//ff:func feature=match type=helper control=selection lang=rust
//ff:what rsCollectCfgTestBody: classifies one direct child of the source_file root and returns the carry-forward #[cfg(test)] pending flag, appending a guarded mod_item's body to out. So only calls inside an in-file #[cfg(test)] module are harvested for content matching — production call sites (a fn calling a helper) are excluded, preventing a never-tested function from being falsely attributed. rsAttrIsCfgTest recognizes the guard. Keeps rsCfgTestBodies a flat loop.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsCollectCfgTestBody appends a #[cfg(test)] mod body to out and returns the
// updated pending flag.
func rsCollectCfgTestBody(c *treesitter.Node, pending bool, out *[]*treesitter.Node) bool {
	switch c.Type {
	case "attribute_item":
		return pending || rsAttrIsCfgTest(c)
	case "mod_item":
		if b := c.ChildByField("body"); pending && b != nil {
			*out = append(*out, b)
		}
		return false
	default:
		return pending
	}
}
