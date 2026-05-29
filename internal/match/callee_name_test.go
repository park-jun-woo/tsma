package match

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// firstCallExpr parses a Go expression statement inside a function and returns
// the first *ast.CallExpr's callee (Fun) found in the body.
func firstCallFun(t *testing.T, src string) ast.Expr {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var found ast.Expr
	ast.Inspect(f, func(n ast.Node) bool {
		if found != nil {
			return false
		}
		if call, ok := n.(*ast.CallExpr); ok {
			found = call.Fun
			return false
		}
		return true
	})
	if found == nil {
		t.Fatalf("no call expr found in %q", src)
	}
	return found
}

func TestCalleeNameIdent(t *testing.T) {
	fun := firstCallFun(t, "package p\nfunc f() { Foo() }\n")
	if got := calleeName(fun); got != "Foo" {
		t.Errorf("calleeName(Ident) = %q, want %q", got, "Foo")
	}
}

func TestCalleeNameSelector(t *testing.T) {
	fun := firstCallFun(t, "package p\nfunc f() { x.Foo() }\n")
	if got := calleeName(fun); got != "Foo" {
		t.Errorf("calleeName(Selector) = %q, want %q", got, "Foo")
	}
}

func TestCalleeNameNestedCall(t *testing.T) {
	// g()() : the outer callee is itself a CallExpr -> "".
	fun := firstCallFun(t, "package p\nfunc f() { g()() }\n")
	if got := calleeName(fun); got != "" {
		t.Errorf("calleeName(CallExpr callee) = %q, want \"\"", got)
	}
}
