//ff:func feature=match type=helper control=sequence lang=java
//ff:what isJavaTestFile: reports whether a filename is a JUnit test file — a .java file whose base ends in Test or Tests (FooTest.java / FooTests.java), the same convention JavaMatcher.Match and isJavaSource (exclusion) use. Drives which files BuildJavaPkgTestIndex ingests.
package match

import "strings"

// isJavaTestFile returns true if name is a *Test.java / *Tests.java file.
func isJavaTestFile(name string) bool {
	if !strings.HasSuffix(name, ".java") {
		return false
	}
	stem := strings.TrimSuffix(name, ".java")
	return strings.HasSuffix(stem, "Test") || strings.HasSuffix(stem, "Tests")
}
