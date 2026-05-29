//ff:func feature=index type=helper control=sequence
//ff:what Appends a Java method declaration to the parse state
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// appendJavaFunc records a method/constructor declaration found on the current
// line. A method is considered exported when declared public.
func appendJavaFunc(st *javaParseState, trimmed, name string) {
	closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
	st.functions = append(st.functions, model.Function{
		QualifiedName: buildJavaQualifiedName(st.pkg, st.scopes, name),
		Name:          name,
		File:          st.relPath,
		StartLine:     st.lineNum,
		EndLine:       st.lineNum,
		Exported:      strings.HasPrefix(trimmed, "public"),
		Status:        model.StatusTodo,
	})
}
