//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what isCsTestFile: reports whether a filename is a C# test file — a .cs file whose base ends in Test or Tests (FooTest.cs / FooTests.cs), the same convention CsMatcher.Match and isCsSource (exclusion) use. Drives which files BuildCsPkgTestIndex ingests.
package match

import "strings"

// isCsTestFile returns true if name is a *Test.cs / *Tests.cs file.
func isCsTestFile(name string) bool {
	if !strings.HasSuffix(name, ".cs") {
		return false
	}
	stem := strings.TrimSuffix(name, ".cs")
	return strings.HasSuffix(stem, "Test") || strings.HasSuffix(stem, "Tests")
}
