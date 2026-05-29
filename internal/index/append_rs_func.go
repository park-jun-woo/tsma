//ff:func feature=index type=helper control=sequence
//ff:what Appends a Rust function declaration to the parse state unless it is a #[cfg(test)] item
package index

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// appendRsFunc records a fn declaration found on the current line, skipping
// functions that belong to a #[cfg(test)] scope.
func appendRsFunc(st *rsParseState, trimmed string) {
	if cfgTestActive(st.scopes, st.pendingCfgTest) {
		return
	}
	m := rsFnPattern.FindStringSubmatch(trimmed)
	closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
	name := m[1]
	st.functions = append(st.functions, model.Function{
		QualifiedName: buildRsQualifiedName(st.relDir, st.scopes, name),
		Name:          name,
		File:          st.relPath,
		StartLine:     st.lineNum,
		EndLine:       st.lineNum,
		Exported:      strings.HasPrefix(trimmed, "pub"),
		Status:        model.StatusTodo,
	})
}
