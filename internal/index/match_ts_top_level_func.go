//ff:func feature=index type=helper control=sequence
//ff:what Matches a top-level TS/JS function declaration and returns a model.Function
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// matchTSTopLevelFunc attempts to match a top-level function declaration.
func matchTSTopLevelFunc(trimmed, relDir, relPath string, lineNum int) (model.Function, bool) {
	m := tsFuncPattern.FindStringSubmatch(trimmed)
	if m == nil {
		return model.Function{}, false
	}

	name := m[1]
	if name == "" {
		name = m[2]
		// Phase007: const/let/var matched but no arrow function or function keyword — skip non-function assignments.
		if name != "" && !strings.Contains(trimmed, "=>") && !strings.Contains(trimmed, "function") {
			return model.Function{}, false
		}
	}
	if name == "" {
		return model.Function{}, false
	}

	qualifiedName := buildQualifiedName(relDir, "", name)
	exported := strings.HasPrefix(trimmed, "export")

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
