package match

import (
	"go/ast"
	"testing"
)

func TestCalleeReceiver(t *testing.T) {
	varTypes := map[string]string{"f": "CSharpFile"}
	cases := []struct {
		src  string
		want string
	}{
		{`(&GoFile{}).GetFuncs()`, "GoFile"}, // composite literal
		{`f.GetFuncs()`, "CSharpFile"},       // local var lookup
		{`g.GetFuncs()`, ""},                 // unknown var
		{`Plain()`, ""},                      // free function (Ident callee)
		{`obj.field.M()`, ""},                // nested selector, X is a SelectorExpr
	}
	for _, c := range cases {
		call, ok := parseExpr(t, c.src).(*ast.CallExpr)
		if !ok {
			t.Fatalf("%q is not a call expr", c.src)
		}
		if got := calleeReceiver(call.Fun, varTypes); got != c.want {
			t.Errorf("calleeReceiver(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}
