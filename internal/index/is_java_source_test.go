package index

import "testing"

func TestIsJavaSource(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"src/main/java/p/Foo.java", true},
		{"Foo.java", true},
		{"src/test/java/p/FooTest.java", false},
		// A non-*Test file under a test source tree (with a leading slash before
		// /src/test/) must be excluded via the path check, not the name check.
		{"project/src/test/java/p/Helper.java", false},
		{"src/main/java/p/FooTest.java", false},
		{"src/main/java/p/FooTests.java", false},
		{"Foo.kt", false},
		{"README.md", false},
	}
	for _, c := range cases {
		if got := isJavaSource(c.path); got != c.want {
			t.Errorf("isJavaSource(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}
