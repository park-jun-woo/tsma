//ff:func feature=match type=helper control=sequence lang=rust
//ff:what rsAttrIsCfgTest: reports whether an attribute_item is a `#[cfg(test)]` guard — its subtree mentions both the `cfg` and `test` identifiers (so `#[cfg(test)]` matches while `#[test]` or `#[cfg(feature="x")]` do not). Lets rsCfgTestBodies isolate the in-file unit-test module for content-aware attribution. Node-based mirror of the index/smell cfg(test) predicates.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsAttrIsCfgTest returns true when attr is a #[cfg(test)] attribute_item.
func rsAttrIsCfgTest(attr *treesitter.Node) bool {
	var hasCfg, hasTest bool
	treesitter.Walk(attr, func(n *treesitter.Node) bool {
		if n.Type != "identifier" {
			return true
		}
		if n.Text == "cfg" {
			hasCfg = true
		}
		if n.Text == "test" {
			hasTest = true
		}
		return true
	})
	return hasCfg && hasTest
}
