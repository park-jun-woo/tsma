//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Builds a Java qualified name from package, enclosing class stack, and method name
package index

import "strings"

// buildJavaQualifiedName constructs a "pkg.Outer.Inner.method" style qualified
// name. The package and each enclosing class/interface/enum from the scope
// stack are joined with ".", then the method name is appended.
func buildJavaQualifiedName(pkg string, scopes []javaScope, name string) string {
	var segs []string
	if pkg != "" {
		segs = append(segs, pkg)
	}
	for _, s := range scopes {
		if s.typeName != "" {
			segs = append(segs, s.typeName)
		}
	}
	segs = append(segs, name)
	return strings.Join(segs, ".")
}
