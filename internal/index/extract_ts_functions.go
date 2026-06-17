//ff:func feature=index type=helper control=sequence lang=typescript
//ff:what extractTSFunctions: the TypeScript grammar interpreter entry for the shared tree-sitter pipeline. Delegates to walkTSTopLevel to emit one model.Function per function_declaration / const-arrow / class method — the exact same model.Function shape the line-based indexTSFile produces, so the matcher and coverage stages are unchanged. 005b~e replace this extractor with their own grammar interpreter.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// extractTSFunctions converts a parsed TypeScript file (tree-sitter root node)
// into the same []model.Function the line-based indexTSFile yields. relDir is
// pkgDirOf(relPath); it scopes qualified names exactly as the fallback does.
func extractTSFunctions(root *treesitter.Node, relPath, relDir string) []model.Function {
	if root == nil {
		return nil
	}
	var out []model.Function
	walkTSTopLevel(root, relDir, relPath, false, &out)
	return out
}
