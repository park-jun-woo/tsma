//ff:func feature=index type=helper control=sequence lang=java
//ff:what extractJavaFunctions: the Java grammar interpreter entry for the shared tree-sitter pipeline. Reads the package declaration from the AST (javaPackageName), then collects one model.Function per method_declaration / constructor_declaration across every (possibly nested) class/interface/enum/record body via collectJavaMembers — the exact same model.Function shape the line-based indexJavaFile produces (buildJavaQualifiedName pkg.Outer.Inner.method), so matcher and coverage stages are unchanged. relDir is unused: Java qualified names are scoped by the declared package, not the directory, matching the fallback.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// extractJavaFunctions converts a parsed Java file (tree-sitter program root)
// into the same []model.Function the line-based indexJavaFile yields. It walks
// only type bodies (never method bodies), matching the indexer's top-level
// semantics, so anonymous classes inside method bodies are not indexed.
func extractJavaFunctions(root *treesitter.Node, relPath, relDir string) []model.Function {
	if root == nil {
		return nil
	}
	pkg := javaPackageName(root)
	var out []model.Function
	collectJavaMembers(root, pkg, nil, relPath, &out)
	return out
}
