package index

import "testing"

func TestPopClosedJavaScopes(t *testing.T) {
	scopes := []javaScope{{depth: 0, typeName: "Outer"}, {depth: 1, typeName: "Inner"}}

	// At depth 2 nothing closes.
	if got := popClosedJavaScopes(scopes, 2); len(got) != 2 {
		t.Errorf("at depth 2 len = %d, want 2", len(got))
	}
	// At depth 1 the inner scope (opened at depth 1) closes.
	got := popClosedJavaScopes(scopes, 1)
	if len(got) != 1 || got[0].typeName != "Outer" {
		t.Errorf("at depth 1 got %+v, want [Outer]", got)
	}
	// At depth 0 both close.
	if got := popClosedJavaScopes(scopes, 0); len(got) != 0 {
		t.Errorf("at depth 0 len = %d, want 0", len(got))
	}
}
