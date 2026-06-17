//ff:func feature=index type=helper control=sequence lang=python
//ff:what pyAstFuncToModel: maps one ast dump entry to a model.Function — the same shape the line-based indexPyFile yields. Builds the qualified name from relDir + receiver (buildQualifiedName, shared with the line path) and reuses the line path's exported heuristic (leading uppercase, no leading underscore). Extracted from parsePyAst's loop so each stays within the per-function line budget.
package index

import (
	"strings"
	"unicode"

	"github.com/park-jun-woo/tsma/internal/model"
)

// pyAstFuncToModel converts a single ast dump entry into a model.Function.
func pyAstFuncToModel(f pyAstFunc, relDir, relPath string) model.Function {
	exported := len(f.Name) > 0 && unicode.IsUpper(rune(f.Name[0])) && !strings.HasPrefix(f.Name, "_")
	return model.Function{
		QualifiedName: buildQualifiedName(relDir, f.Receiver, f.Name),
		Name:          f.Name,
		Receiver:      f.Receiver,
		File:          relPath,
		StartLine:     f.StartLine,
		EndLine:       f.EndLine,
		Exported:      exported,
		Status:        model.StatusTodo,
	}
}
