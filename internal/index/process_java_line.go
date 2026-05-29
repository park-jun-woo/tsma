//ff:func feature=index type=helper control=sequence
//ff:what Processes one Java source line: records package, dispatches declarations, and updates brace depth/scopes
package index

import "strings"

// processJavaLine advances the parse state by one source line. It records the
// last non-empty line, captures the package declaration, skips annotation
// lines, dispatches type/method declarations, then updates brace depth and
// pops closed scopes.
func processJavaLine(st *javaParseState, line string) {
	st.lineNum++
	trimmed := strings.TrimSpace(line)

	if trimmed != "" {
		st.lastNonEmptyLine = st.lineNum
	}

	if m := javaPackagePattern.FindStringSubmatch(trimmed); m != nil {
		st.pkg = m[1]
	}

	if !isJavaAnnotation(trimmed) {
		dispatchJavaLine(st, trimmed)
	}

	st.depth += countBraces(line)
	st.scopes = popClosedJavaScopes(st.scopes, st.depth)
}
