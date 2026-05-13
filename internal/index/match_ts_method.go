//ff:func feature=index type=helper control=sequence
//ff:what Matches a TS/JS class method declaration and returns a model.Function
package index

import (
	"unicode"

	"github.com/park-jun-woo/tsma/internal/model"
)

// tsMethodSkipNames are keywords that look like method calls but should be skipped.
var tsMethodSkipNames = map[string]bool{
	"constructor": true, "if": true, "for": true, "while": true, "switch": true,
}

// matchTSMethod attempts to match a class method declaration.
func matchTSMethod(line, currentClass, relDir, relPath string, lineNum int) (model.Function, bool) {
	m := tsMethodPattern.FindStringSubmatch(line)
	if m == nil {
		return model.Function{}, false
	}

	name := m[1]
	if tsMethodSkipNames[name] {
		return model.Function{}, false
	}

	qualifiedName := buildQualifiedName(relDir, currentClass, name)
	exported := unicode.IsUpper(rune(name[0]))

	return model.Function{
		QualifiedName: qualifiedName,
		Name:          name,
		File:          relPath,
		StartLine:     lineNum,
		EndLine:       lineNum,
		Exported:      exported,
		Status:        model.StatusTodo,
	}, true
}
