//ff:func feature=index type=helper control=sequence
//ff:what Dispatches a Java source line to class-scope push or method declaration handling
package index

import "strings"

// dispatchJavaLine classifies the current line as a type declaration (pushing a
// class/interface/enum scope) or a method declaration, updating the parse state.
func dispatchJavaLine(st *javaParseState, trimmed string) {
	if m := javaTypePattern.FindStringSubmatch(trimmed); m != nil && strings.Contains(trimmed, "{") {
		closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
		st.scopes = append(st.scopes, javaScope{depth: st.depth, typeName: m[1]})
		return
	}
	if name, ok := matchJavaMethod(trimmed); ok {
		appendJavaFunc(st, trimmed, name)
	}
}
