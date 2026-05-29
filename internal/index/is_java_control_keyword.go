//ff:func feature=index type=helper control=selection
//ff:what Returns true if a captured Java identifier is a control-flow keyword rather than a method name
package index

// isJavaControlKeyword reports whether name is a Java control-flow or statement
// keyword that can syntactically look like a method call followed by a brace
// (e.g. `if (...) {`, `for (...) {`). Such lines must not be treated as method
// declarations by the indexer.
func isJavaControlKeyword(name string) bool {
	switch name {
	case "if", "for", "while", "switch", "catch", "synchronized",
		"return", "else", "do", "try", "finally", "new", "case":
		return true
	}
	return false
}
