//ff:func feature=index type=helper control=sequence
//ff:what Processes one Rust source line: tracks attributes, dispatches declarations, and updates brace depth/scopes
package index

import "strings"

// processRsLine advances the parse state by one source line. It records the
// last non-empty line, handles #[cfg(test)] attributes, dispatches impl/mod/fn
// declarations, then updates brace depth and pops closed scopes.
func processRsLine(st *rsParseState, line string) {
	st.lineNum++
	trimmed := strings.TrimSpace(line)

	if trimmed != "" {
		st.lastNonEmptyLine = st.lineNum
	}

	if handleRsAttribute(st, trimmed) {
		return
	}

	dispatchRsLine(st, trimmed, line)

	st.depth += countBraces(line)
	st.scopes = popClosedScopes(st.scopes, st.depth)
	st.pendingCfgTest = false
}
