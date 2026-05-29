//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Dispatches a C# source line to namespace/type scope handling or method declaration
package index

import "strings"

// dispatchCsLine classifies the current line as a namespace declaration, a type
// declaration (class/struct/interface/record/enum), or a method declaration,
// updating the parse state.
//
// A file-scoped namespace ("namespace X;", no opening brace) is recorded in
// fileNs immediately. Block-scoped namespaces and types are appended to the
// pending list; processCsLine flushes pending names onto the scope stack when
// the opening brace is seen, which supports both K&R (brace on the declaration
// line) and Allman (brace on the next line) styles.
func dispatchCsLine(st *csParseState, trimmed string) {
	if m := csNamespacePattern.FindStringSubmatch(trimmed); m != nil {
		if strings.Contains(trimmed, ";") && !strings.Contains(trimmed, "{") {
			st.fileNs = m[1]
		} else {
			st.pending = append(st.pending, m[1])
		}
		return
	}
	if m := csTypePattern.FindStringSubmatch(trimmed); m != nil {
		closePrevTSEndLine(st.functions, st.lineNum, st.lastNonEmptyLine)
		st.pending = append(st.pending, m[1])
		return
	}
	if name, ok := matchCsMethod(trimmed); ok {
		appendCsFunc(st, trimmed, name)
	}
}
