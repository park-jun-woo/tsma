//ff:func feature=index type=helper control=sequence
//ff:what Matches a Python function definition and returns a model.Function with class context
package index

import (
	"strings"
	"unicode"

	"github.com/park-jun-woo/tsma/internal/model"
)

// matchPyFunc attempts to match a Python function definition.
// Returns the function (or nil) and the updated currentClass value.
func matchPyFunc(line string, lineNum int, relPath, relDir, currentClass string, classIndent int) (*model.Function, string) {
	m := pyFuncPattern.FindStringSubmatch(line)
	if m == nil {
		return nil, currentClass
	}

	indent := pyIndent(m[1])
	name := m[2]

	// Reset class context if indent is at or below class level.
	if currentClass != "" && indent <= classIndent {
		currentClass = ""
	}

	isMethod := currentClass != "" && indent > classIndent

	receiver := ""
	if isMethod {
		receiver = currentClass
	}

	qualifiedName := buildQualifiedName(relDir, receiver, name)
	exported := len(name) > 0 && unicode.IsUpper(rune(name[0])) && !strings.HasPrefix(name, "_")

	fn := &model.Function{
		QualifiedName: qualifiedName,
		Name:          name,
		File:          relPath,
		StartLine:     lineNum,
		EndLine:       lineNum,
		Exported:      exported,
		Status:        model.StatusTodo,
	}

	return fn, currentClass
}
