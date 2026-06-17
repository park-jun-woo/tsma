//ff:func feature=index type=helper control=sequence lang=typescript
//ff:what tsMethodFromDefinition: builds a model.Function from a method_definition node, embedding the receiver class in the qualified name (Receiver left empty to match the line-based output). `constructor` is skipped; exported mirrors the existing uppercase-first heuristic so precise path and fallback agree.
package index

import (
	"unicode"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// tsSkipMethodNames are method names never indexed as testable functions. Only
// `constructor` is a real method_definition node (if/for/while/switch are
// statements, so the AST already excludes them — unlike the line-based regex).
var tsSkipMethodNames = map[string]bool{
	"constructor": true,
}

// tsMethodFromDefinition converts a method_definition into a model.Function, or
// returns false for a nameless or skipped method.
func tsMethodFromDefinition(m *treesitter.Node, className, relDir, relPath string) (model.Function, bool) {
	if m.Type != "method_definition" {
		return model.Function{}, false
	}
	nameNode := m.ChildByField("name")
	if nameNode == nil || nameNode.Text == "" {
		return model.Function{}, false
	}
	name := nameNode.Text
	if tsSkipMethodNames[name] {
		return model.Function{}, false
	}
	return model.Function{
		QualifiedName: buildQualifiedName(relDir, className, name),
		Name:          name,
		File:          relPath,
		StartLine:     m.StartLine(),
		EndLine:       m.EndLine(),
		Exported:      unicode.IsUpper(rune(name[0])),
		Status:        model.StatusTodo,
	}, true
}
