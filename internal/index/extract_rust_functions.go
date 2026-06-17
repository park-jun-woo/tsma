//ff:func feature=index type=helper control=sequence lang=rust
//ff:what extractRustFunctions: the Rust grammar interpreter entry for the shared tree-sitter pipeline. Collects one model.Function per free function_item and impl method across every (possibly nested) mod/impl body via collectRsMembers — the exact same model.Function shape the line-based indexRsFile produces (buildRsQualifiedName relDir + mod::Recv::name), so matcher and coverage stages are unchanged. Functions inside a #[cfg(test)] mod (the in-file unit-test module) are excluded, matching appendRsFunc's cfgTestActive skip. relDir scopes the qualified name exactly as the fallback's pkgDirOf does.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// extractRustFunctions converts a parsed Rust file (tree-sitter source_file root)
// into the same []model.Function the line-based indexRsFile yields. It walks only
// mod/impl bodies (never function bodies), so local functions inside a function
// body are not indexed, matching the line-based top-level semantics.
func extractRustFunctions(root *treesitter.Node, relPath, relDir string) []model.Function {
	if root == nil {
		return nil
	}
	var out []model.Function
	collectRsMembers(root, relDir, nil, relPath, &out)
	return out
}
