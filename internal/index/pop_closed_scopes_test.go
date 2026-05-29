package index

import "testing"

func TestPopClosedScopes(t *testing.T) {
	scopes := []rsScope{{depth: 0, receiver: "Foo"}, {depth: 1, module: "inner"}}

	// Still deep inside both scopes.
	if got := popClosedScopes(scopes, 2); len(got) != 2 {
		t.Errorf("depth 2: kept %d scopes, want 2", len(got))
	}
	// Closed the inner scope (depth back to 1).
	if got := popClosedScopes(scopes, 1); len(got) != 1 {
		t.Errorf("depth 1: kept %d scopes, want 1", len(got))
	}
	// Closed everything.
	if got := popClosedScopes(scopes, 0); len(got) != 0 {
		t.Errorf("depth 0: kept %d scopes, want 0", len(got))
	}
}
