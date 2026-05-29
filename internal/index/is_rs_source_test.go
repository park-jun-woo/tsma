package index

import "testing"

func TestIsRsSource(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"src/lib.rs", true},
		{"src/foo/bar.rs", true},
		{"main.rs", true},
		{"build.rs", false},
		{"src/build.rs", false},
		{"README.md", false},
		{"src/lib.go", false},
	}
	for _, c := range cases {
		if got := isRsSource(c.path); got != c.want {
			t.Errorf("isRsSource(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}
