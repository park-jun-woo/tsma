//ff:func feature=index type=helper control=sequence
//ff:what Returns true if the path is a non-test Java source file eligible for indexing
package index

import "strings"

// isJavaSource returns true if the path is a Java source file eligible for
// indexing. Files under a test source tree (src/test/) and *Test.java /
// *Tests.java files are excluded, since indexing targets production code.
func isJavaSource(path string) bool {
	if !strings.HasSuffix(path, ".java") {
		return false
	}
	slashed := strings.ReplaceAll(path, "\\", "/")
	if strings.Contains(slashed, "/src/test/") {
		return false
	}
	base := slashed
	if idx := strings.LastIndex(slashed, "/"); idx >= 0 {
		base = slashed[idx+1:]
	}
	name := strings.TrimSuffix(base, ".java")
	if strings.HasSuffix(name, "Test") || strings.HasSuffix(name, "Tests") {
		return false
	}
	return true
}
