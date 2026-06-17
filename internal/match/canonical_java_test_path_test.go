package match

import (
	"path/filepath"
	"testing"
)

// TestCanonicalJavaTestPath asserts the JUnit mirror formula: a src/main source
// maps to the FooTest.java sibling under src/test (the javaTestDir SSOT), and a
// source outside a src/main tree falls back to the same directory.
func TestCanonicalJavaTestPath(t *testing.T) {
	cases := []struct {
		lang   string
		source string
		want   string
	}{
		{"java", "src/main/java/com/example/calc/Calculator.java",
			"src/test/java/com/example/calc/CalculatorTest.java"},
		{"java", "lib/Foo.java", "lib/FooTest.java"},
		{"java", "src/main/java/Top.java", "src/test/java/TopTest.java"},
		// Non-.java and non-java language must return "".
		{"java", "src/main/java/com/example/README.md", ""},
		{"go", "pkg/foo.java", ""},
	}
	for _, c := range cases {
		got := CanonicalTestPath(c.lang, c.source)
		if filepath.ToSlash(got) != c.want {
			t.Errorf("CanonicalTestPath(%q, %q) = %q, want %q", c.lang, c.source, got, c.want)
		}
	}
}
