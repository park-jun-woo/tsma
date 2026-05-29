package index

import "testing"

func TestFlushPendingCsScopes(t *testing.T) {
	scopes := []csScope{{depth: 0, typeName: "Ns"}}
	scopes, pending := flushPendingCsScopes(scopes, []string{"Outer"}, 1)
	if len(pending) != 0 {
		t.Errorf("pending should be cleared, got %v", pending)
	}
	if len(scopes) != 2 {
		t.Fatalf("scopes = %d, want 2", len(scopes))
	}
	if scopes[1].typeName != "Outer" || scopes[1].depth != 1 {
		t.Errorf("flushed scope = %+v, want {depth:1 Outer}", scopes[1])
	}
}

func TestFlushPendingCsScopesEmpty(t *testing.T) {
	scopes := []csScope{{depth: 0, typeName: "Ns"}}
	got, pending := flushPendingCsScopes(scopes, nil, 1)
	if len(got) != 1 || pending != nil {
		t.Errorf("no pending: scopes=%+v pending=%v", got, pending)
	}
}
