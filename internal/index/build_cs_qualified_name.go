//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what Builds a C# qualified name from file-scoped namespace, enclosing scope stack, and method name
package index

import "strings"

// buildCsQualifiedName constructs a "Namespace.Outer.Inner.Method" style
// qualified name. A file-scoped namespace (fileNs) and each enclosing
// namespace/type from the scope stack are joined with ".", then the method name
// is appended.
func buildCsQualifiedName(fileNs string, scopes []csScope, name string) string {
	var segs []string
	if fileNs != "" {
		segs = append(segs, fileNs)
	}
	for _, s := range scopes {
		if s.typeName != "" {
			segs = append(segs, s.typeName)
		}
	}
	segs = append(segs, name)
	return strings.Join(segs, ".")
}
