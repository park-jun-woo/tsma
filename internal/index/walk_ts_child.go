//ff:func feature=index type=helper control=selection dimension=1 lang=typescript
//ff:what walkTSChild: dispatches one top-level child by node type — unwraps export_statement (exported=true) back into walkTSTopLevel, emits a function_declaration, or collects const-arrow declarators / class methods. The single place the three indexable TS shapes are recognized.
package index

import (
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// walkTSChild handles one child node at top-level scope.
func walkTSChild(c *treesitter.Node, relDir, relPath string, exported bool, out *[]model.Function) {
	switch c.Type {
	case "export_statement":
		walkTSTopLevel(c, relDir, relPath, true, out)
	case "function_declaration", "generator_function_declaration":
		if fn, ok := tsFuncFromDecl(c, relDir, relPath, exported); ok {
			*out = append(*out, fn)
		}
	case "lexical_declaration", "variable_declaration":
		collectTSArrowFuncs(c, relDir, relPath, exported, out)
	case "class_declaration", "abstract_class_declaration":
		collectTSMethods(c, relDir, relPath, out)
	}
}
