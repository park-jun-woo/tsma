package match

import "testing"

func TestJavaTestDir(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"src/main/java/com/example", "src/test/java/com/example"},
		{"src/main/java", "src/test/java"},
		{"lib", "lib"}, // no main tree -> beside source
	}
	for _, c := range cases {
		if got := javaTestDir(c.src); got != c.want {
			t.Errorf("javaTestDir(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}
