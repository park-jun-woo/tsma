//ff:func feature=index type=helper control=selection dimension=1 lang=rust
//ff:what dispatchRsMember: classifies one member node of a Rust source_file / mod body / impl body and returns the carry-forward #[cfg(test)] pending flag. An attribute_item folds its cfg(test)-ness into pending (so the next item is excluded); a function_item is emitted as a model.Function unless cfgTestActive (the in-file #[cfg(test)] mod exclusion); an impl_item / mod_item recurses via walkRsImpl / walkRsMod, propagating the resolved cfgTest so nested functions inherit the test guard. Any other node (use, comment, brace token) preserves pending unchanged. The single place Rust's indexable shapes are recognized (the dispatchCSMember analogue).
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// dispatchRsMember handles one member node by type and returns the updated
// pending #[cfg(test)] flag for the next sibling.
func dispatchRsMember(c *treesitter.Node, relDir string, scopes []rsScope, relPath string, out *[]model.Function, pending bool) bool {
	switch c.Type {
	case "attribute_item":
		return pending || rsAttrCfgTest(c)
	case "function_item":
		emitRsFunc(c, relDir, scopes, relPath, out, pending)
		return false
	case "impl_item":
		walkRsImpl(c, relDir, scopes, relPath, out, cfgTestActive(scopes, pending))
		return false
	case "mod_item":
		walkRsMod(c, relDir, scopes, relPath, out, cfgTestActive(scopes, pending))
		return false
	default:
		return pending
	}
}
