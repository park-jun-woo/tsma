//ff:func feature=runner type=helper control=sequence lang=csharp
//ff:what Derives the test class name from a C# test file path
package runner

import (
	"path/filepath"
	"strings"
)

// csTestClass returns the simple test class name for a C# test file path,
// e.g. "App.Tests/Services/FooTests.cs" -> "FooTests". The .NET test filter
// matches the fully-qualified name by substring, so the simple class name is
// sufficient to select a single test class.
func csTestClass(testFile string) string {
	return strings.TrimSuffix(filepath.Base(testFile), ".cs")
}
