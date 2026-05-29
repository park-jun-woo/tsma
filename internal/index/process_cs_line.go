//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Processes one C# source line: dispatches declarations, flushes pending scopes, and updates brace depth
package index

import "strings"

// processCsLine advances the parse state by one source line. It records the last
// non-empty line, skips attribute lines, dispatches namespace/type/method
// declarations, flushes any pending namespace/type scope at the opening brace
// (so K&R and Allman brace styles both work), then updates brace depth and pops
// closed scopes.
func processCsLine(st *csParseState, line string) {
	st.lineNum++
	trimmed := strings.TrimSpace(line)

	if trimmed != "" {
		st.lastNonEmptyLine = st.lineNum
	}

	if !isCsAttribute(trimmed) {
		dispatchCsLine(st, trimmed)
	}

	if strings.Contains(line, "{") {
		st.scopes, st.pending = flushPendingCsScopes(st.scopes, st.pending, st.depth)
	}

	st.depth += countBraces(line)
	st.scopes = popClosedCsScopes(st.scopes, st.depth)
}
