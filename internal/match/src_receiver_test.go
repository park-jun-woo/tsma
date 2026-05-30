package match

import "testing"

func TestSrcReceiver(t *testing.T) {
	cases := []struct {
		src  string // a receiver type expression
		want string
	}{
		{`T`, "T"},
		{`*T`, "T"},
		{`T[X]`, "T"},
		{`T[X, Y]`, "T"},
		{`*T[X]`, "T"},
		{`pkg.T`, ""}, // selector receiver is not a valid Go method receiver -> ""
	}
	for _, c := range cases {
		if got := srcReceiver(parseExpr(t, c.src)); got != c.want {
			t.Errorf("srcReceiver(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}
