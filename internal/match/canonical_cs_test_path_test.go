package match

import (
	"path/filepath"
	"testing"
)

// TestCanonicalCsTestPath asserts the C# source→test-project naming formula:
// Calc/Foo.cs maps to the parallel Calc.Tests/FooTests.cs (xUnit), a nested
// source keeps its sub-path under the *.Tests project, and a non-.cs file or a
// project-less source degrades predictably.
func TestCanonicalCsTestPath(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"Calc/Calculator.cs", "Calc.Tests/Calculator" + "Tests.cs"},
		{"App/Services/Foo.cs", "App.Tests/Services/FooTests.cs"},
		{"Foo.cs", "FooTests.cs"},
		{"Calc/Calculator.txt", ""},
	}
	for _, c := range cases {
		got := filepath.ToSlash(CanonicalTestPath("csharp", c.src))
		if got != c.want {
			t.Errorf("CanonicalTestPath(csharp, %q) = %q, want %q", c.src, got, c.want)
		}
	}
}
