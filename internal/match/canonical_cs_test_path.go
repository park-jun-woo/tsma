//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what canonicalCsTestPath: the C# arm of CanonicalTestPath (xUnit/NUnit) — Calc/Foo.cs → the parallel test project Calc.Tests/FooTests.cs. The source→test-project directory transform is delegated to csTestDirs (the one SSOT also used by CsMatcher.Match), choosing its last candidate (the *.Tests project dir) so the loop's write path and the matcher's read path agree. Non-.cs files return "".
package match

import (
	"path/filepath"
	"strings"
)

// canonicalCsTestPath maps a C# source file to its test-project path
// (FooTests.cs under the parallel *.Tests project), or returns "" when base is
// not .cs.
func canonicalCsTestPath(sourceFile, base string) string {
	if !strings.HasSuffix(base, ".cs") {
		return ""
	}
	dirs := csTestDirs(filepath.Dir(sourceFile))
	testDir := dirs[len(dirs)-1]
	testBase := strings.TrimSuffix(base, ".cs") + "Tests.cs"
	return filepath.Join(testDir, testBase)
}
