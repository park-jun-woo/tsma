//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what extractCSharpFunctions: the C# grammar interpreter entry for the shared tree-sitter pipeline. Reads the file-scoped namespace from the AST (csFileNamespace), then collects one model.Function per method_declaration / constructor_declaration / property_declaration across every (possibly nested) namespace/class/struct/interface/record/enum body via collectCSMembers — the exact same model.Function shape the line-based indexCsFile produces (buildCsQualifiedName Namespace.Outer.Inner.Member), so matcher and coverage stages are unchanged. relDir is unused: C# qualified names are scoped by the declared namespace, not the directory, matching the fallback.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// extractCSharpFunctions converts a parsed C# file (tree-sitter compilation_unit
// root) into the same []model.Function the line-based indexCsFile yields. It
// walks only type/namespace bodies (never method bodies), matching the indexer's
// top-level semantics, so local functions inside method bodies are not indexed.
func extractCSharpFunctions(root *treesitter.Node, relPath, relDir string) []model.Function {
	if root == nil {
		return nil
	}
	fileNs := csFileNamespace(root)
	var out []model.Function
	collectCSMembers(root, fileNs, nil, relPath, &out)
	return out
}
