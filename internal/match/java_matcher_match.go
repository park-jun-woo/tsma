//ff:func feature=match type=implementation control=iteration dimension=1
//ff:what Looks for FooTest.java / FooTests.java under the Maven/Gradle test source tree
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// Match finds the test file covering the given Java source file.
//
// For the standard Maven/Gradle layout, src/main/java/p/Foo.java maps to
// src/test/java/p/FooTest.java (or FooTests.java). When the source is not under
// a src/main tree, the same directory is searched. The returned testFile is a
// projectRoot-relative path; (\"\", false) is returned when no test exists.
func (m *JavaMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	srcDir := filepath.Dir(sourceFile)
	base := strings.TrimSuffix(filepath.Base(sourceFile), ".java")
	testDir := javaTestDir(srcDir)

	for _, suffix := range []string{"Test", "Tests"} {
		testRel := filepath.Join(testDir, base+suffix+".java")
		if _, err := os.Stat(filepath.Join(projectRoot, testRel)); err == nil {
			return testRel, true
		}
	}

	return "", false
}
