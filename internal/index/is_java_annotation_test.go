package index

import "testing"

func TestIsJavaAnnotation(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"@Override", true},
		{`@SuppressWarnings("unchecked")`, true},
		{"@Test", true},
		{"public void foo() {", false},
		{"class Foo {", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isJavaAnnotation(c.line); got != c.want {
			t.Errorf("isJavaAnnotation(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}
