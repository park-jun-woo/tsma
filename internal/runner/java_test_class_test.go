package runner

import (
	"path/filepath"
	"testing"
)

func TestJavaTestClass(t *testing.T) {
	cases := map[string]string{
		filepath.Join("src", "test", "java", "p", "FooTest.java"): "FooTest",
		filepath.Join("src", "test", "java", "BarTests.java"):     "BarTests",
		"Baz.java": "Baz",
	}
	for in, want := range cases {
		if got := javaTestClass(in); got != want {
			t.Errorf("javaTestClass(%q) = %q, want %q", in, got, want)
		}
	}
}
