//ff:func feature=index type=helper control=selection lang=csharp
//ff:what Returns true if a captured C# identifier is a control-flow keyword rather than a method name
package index

// isCsControlKeyword reports whether name is a C# control-flow or statement
// keyword that can syntactically look like a method call followed by a brace
// (e.g. `if (...) {`, `foreach (...) {`). Such lines must not be treated as
// method declarations by the indexer.
func isCsControlKeyword(name string) bool {
	switch name {
	case "if", "for", "foreach", "while", "switch", "catch", "lock",
		"using", "fixed", "return", "else", "do", "try", "finally",
		"new", "case", "get", "set", "add", "remove":
		return true
	}
	return false
}
