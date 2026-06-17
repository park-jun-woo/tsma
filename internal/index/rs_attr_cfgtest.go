//ff:func feature=index type=helper control=sequence lang=rust
//ff:what rsAttrCfgTest: reports whether an attribute_item node is a `#[cfg(test)]` guard — the node-based analogue of the line-based handleRsAttribute's `strings.Contains(trimmed, "cfg(test)")`. It walks the attribute subtree and is satisfied only when both the `cfg` identifier and a `test` identifier appear (so `#[cfg(test)]` and `#[cfg(all(test, feature="x"))]` match, while `#[test]` or `#[cfg(feature="x")]` do not). Drives the same cfgTest exclusion the line-based path applies via cfgTestActive.
package index

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsAttrCfgTest returns true when attr is a #[cfg(test)] attribute_item, i.e. its
// subtree mentions both the `cfg` and `test` identifiers.
func rsAttrCfgTest(attr *treesitter.Node) bool {
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
