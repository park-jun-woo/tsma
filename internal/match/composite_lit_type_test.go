package match

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// parseExpr parses a Go expression and returns its AST.
func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	e, err := parser.ParseExprFrom(token.NewFileSet(), "x.go", src, 0)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	return e
}

func TestCompositeLitType(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{`T{}`, "T"},
		{`&T{}`, "T"},
		{`(&T{})`, "T"},
		{`T{a: 1}`, "T"},
		{`pkg.T{}`, ""}, // selector type, not a bare Ident
		{`x`, ""},       // bare variable
		{`NewT()`, ""},  // constructor call
		{`*p`, ""},      // dereference (StarExpr), not &
		{`[]int{}`, ""}, // composite literal but not an Ident type
		{`-x`, ""},      // UnaryExpr whose op is not & (negation)
		{`<-ch`, ""},    // UnaryExpr whose op is not & (channel receive)
	}
	for _, c := range cases {
		if got := compositeLitType(parseExpr(t, c.src)); got != c.want {
			t.Errorf("compositeLitType(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}
