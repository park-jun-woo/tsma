package index

import "testing"

func TestSkipQuoted(t *testing.T) {
	cases := []struct {
		line  string
		start int
		want  int
	}{
		{`"abc"`, 0, 4},               // string closes at index 4
		{`"a\"b"`, 0, 5},              // escaped quote inside string
		{`'a'`, 0, 2},                 // char literal
		{`'\n'`, 0, 3},                // escaped char literal
		{`&'a str`, 1, 1},             // lifetime: returned unchanged
		{`"unterminated`, 0, 12},      // unterminated string -> last index
		{`'ab`, 0, 0},                 // unterminated single quote (short) -> start
		{`'a`, 0, 0},                  // single quote, one trailing char, EOL -> start
	}
	for _, c := range cases {
		if got := skipQuoted(c.line, c.start); got != c.want {
			t.Errorf("skipQuoted(%q, %d) = %d, want %d", c.line, c.start, got, c.want)
		}
	}
}
