//ff:func feature=runner type=helper control=sequence
//ff:what Derives the JUnit test class name from a Java test file path
package runner

import (
	"path/filepath"
	"strings"
)

// javaTestClass returns the simple test class name for a Java test file path,
// e.g. "src/test/java/p/FooTest.java" -> "FooTest". Build tools accept the
// simple class name for single-class test selection (mvn -Dtest, gradle
// --tests).
func javaTestClass(testFile string) string {
	return strings.TrimSuffix(filepath.Base(testFile), ".java")
}
