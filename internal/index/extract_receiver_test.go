package index

import (
	"go/ast"
	"testing"
)

func TestExtractReceiverIdent(t *testing.T) {
	expr := &ast.Ident{Name: "Handler"}
	got := extractReceiver(expr)
	if got != "Handler" {
		t.Errorf("extractReceiver(Ident) = %q, want %q", got, "Handler")
	}
}

func TestExtractReceiverStarExpr(t *testing.T) {
	expr := &ast.StarExpr{X: &ast.Ident{Name: "Handler"}}
	got := extractReceiver(expr)
	if got != "Handler" {
		t.Errorf("extractReceiver(StarExpr) = %q, want %q", got, "Handler")
	}
}

func TestExtractReceiverIndexExpr(t *testing.T) {
	expr := &ast.IndexExpr{X: &ast.Ident{Name: "Container"}}
	got := extractReceiver(expr)
	if got != "Container" {
		t.Errorf("extractReceiver(IndexExpr) = %q, want %q", got, "Container")
	}
}

func TestExtractReceiverIndexListExpr(t *testing.T) {
	expr := &ast.IndexListExpr{X: &ast.Ident{Name: "MultiGeneric"}}
	got := extractReceiver(expr)
	if got != "MultiGeneric" {
		t.Errorf("extractReceiver(IndexListExpr) = %q, want %q", got, "MultiGeneric")
	}
}

func TestExtractReceiverNestedStarIndex(t *testing.T) {
	// *Container[T]
	expr := &ast.StarExpr{X: &ast.IndexExpr{X: &ast.Ident{Name: "Container"}}}
	got := extractReceiver(expr)
	if got != "Container" {
		t.Errorf("extractReceiver(nested Star+Index) = %q, want %q", got, "Container")
	}
}

func TestExtractReceiverUnknown(t *testing.T) {
	// Use an expression type not handled (e.g. SelectorExpr)
	expr := &ast.SelectorExpr{
		X:   &ast.Ident{Name: "pkg"},
		Sel: &ast.Ident{Name: "Type"},
	}
	got := extractReceiver(expr)
	if got != "Unknown" {
		t.Errorf("extractReceiver(SelectorExpr) = %q, want %q", got, "Unknown")
	}
}
