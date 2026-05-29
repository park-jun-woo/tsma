//ff:func feature=index type=helper control=sequence lang=csharp
//ff:what Returns true if a trimmed C# line is an attribute that should be skipped
package index

import "strings"

// isCsAttribute reports whether the trimmed line is a C# attribute such as
// `[Fact]`, `[TestMethod]`, or `[Obsolete("x")]`. Attribute lines precede a
// declaration and must not be parsed as method/type declarations themselves.
func isCsAttribute(trimmed string) bool {
	return strings.HasPrefix(trimmed, "[")
}
