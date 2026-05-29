//ff:func feature=match type=implementation control=iteration dimension=2 lang=csharp
//ff:what Looks for FooTests.cs / FooTest.cs in the same dir and the parallel *.Tests project dir
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// Match finds the test file covering the given C# source file.
//
// Following the common .NET convention of a separate test project, Foo.cs maps
// to FooTests.cs (preferred) or FooTest.cs, searched first in the same
// directory and then in the parallel *.Tests project directory (e.g.
// App/Services/Foo.cs -> App.Tests/Services/FooTests.cs). The returned testFile
// is a projectRoot-relative path; ("", false) is returned when no test exists.
func (m *CsMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	srcDir := filepath.Dir(sourceFile)
	base := strings.TrimSuffix(filepath.Base(sourceFile), ".cs")

	for _, dir := range csTestDirs(srcDir) {
		for _, suffix := range []string{"Tests", "Test"} {
			testRel := filepath.Join(dir, base+suffix+".cs")
			if _, err := os.Stat(filepath.Join(projectRoot, testRel)); err == nil {
				return testRel, true
			}
		}
	}

	return "", false
}
