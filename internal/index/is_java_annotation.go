//ff:func feature=index type=helper control=sequence
//ff:what Returns true if a trimmed Java line is an annotation that should be skipped
package index

import "strings"

// isJavaAnnotation reports whether the trimmed line is a Java annotation such as
// `@Override`, `@Test`, or `@SuppressWarnings("x")`. Annotation lines precede a
// declaration and must not be parsed as method/type declarations themselves.
func isJavaAnnotation(trimmed string) bool {
	return strings.HasPrefix(trimmed, "@")
}
