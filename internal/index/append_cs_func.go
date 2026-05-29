//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Appends a C# method declaration to the parse state
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// appendCsFunc records a method/constructor declaration found on the current
// line. A method is considered exported when declared public.
func appendCsFunc(st *csParseState, trimmed, name string) {
	closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
	st.functions = append(st.functions, model.Function{
		QualifiedName: buildCsQualifiedName(st.fileNs, st.scopes, name),
		Name:          name,
		File:          st.relPath,
		StartLine:     st.lineNum,
		EndLine:       st.lineNum,
		Exported:      strings.HasPrefix(trimmed, "public"),
		Status:        model.StatusTodo,
	})
}
