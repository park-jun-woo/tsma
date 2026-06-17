//ff:func feature=smell type=helper control=sequence lang=rust
//ff:what rsAttrMentionsTest: reports whether an attribute_item subtree mentions the `test` identifier — true for both `#[test]` and `#[cfg(test)]`. Used by rsTestScopeNodes to decide which mod/function bodies are test scopes, so the escape-hatch detectors only run on test code (the plan's "테스트코드 한정", false-positive zero for legitimate production unsafe/FFI).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsAttrMentionsTest returns true when the attribute_item subtree contains a
// `test` identifier (covering both #[test] and #[cfg(test)]).
func rsAttrMentionsTest(attr *treesitter.Node) bool {
	found := false
	treesitter.Walk(attr, func(n *treesitter.Node) bool {
		if n.Type == "identifier" && n.Text == "test" {
			found = true
		}
		return true
	})
	return found
}
