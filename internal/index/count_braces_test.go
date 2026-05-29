package index

import "testing"

func TestCountBraces(t *testing.T) {
	cases := []struct {
		line string
		want int
	}{
		{"fn foo() {", 1},
		{"}", -1},
		{"fn one() {}", 0},
		{"if x { y } else { z }", 0},
		{`let s = "a { b }";`, 0},
		{`let c = '{';`, 0},
		{"foo(); // trailing { comment", 0},
		{`let s = "with \" quote {";`, 0},
		{`let s = "x"; let t = "{";`, 0},
		{"   ", 0},
		{"fn f<'a>(x: &'a str) {", 1},
		{"let c = '\\n';", 0},
		{"impl<'a> Trait for Foo<'a> {", 1},
	}
	for _, c := range cases {
		if got := countBraces(c.line); got != c.want {
			t.Errorf("countBraces(%q) = %d, want %d", c.line, got, c.want)
		}
	}
}
