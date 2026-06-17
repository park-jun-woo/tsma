//ff:func feature=match type=helper control=sequence lang=java
//ff:what canonicalJavaTestPath: the Java arm of CanonicalTestPath (JUnit) — src/main/java/p/Foo.java → src/test/java/p/FooTest.java in the mirror test tree. The src/main→src/test transform is delegated to javaTestDir (the one SSOT also used by JavaMatcher.Match), so the loop's write path and the matcher's read path agree. Non-.java files return "".
package match

import (
	"path/filepath"
	"strings"
)

// canonicalJavaTestPath maps a Java source file to its JUnit mirror test path
// (FooTest.java under the src/test tree), or returns "" when base is not .java.
func canonicalJavaTestPath(sourceFile, base string) string {
	if !strings.HasSuffix(base, ".java") {
		return ""
	}
	testDir := javaTestDir(filepath.Dir(sourceFile))
	testBase := strings.TrimSuffix(base, ".java") + "Test.java"
	return filepath.Join(testDir, testBase)
}
