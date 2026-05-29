package index

import "testing"

func TestPopClosedCsScopes(t *testing.T) {
	scopes := []csScope{
		{depth: 0, typeName: "Ns"},
		{depth: 1, typeName: "Outer"},
		{depth: 2, typeName: "Inner"},
	}

	// At depth 3 nothing closes.
	got := popClosedCsScopes(scopes, 3)
	if len(got) != 3 {
		t.Fatalf("depth 3: got %d scopes, want 3", len(got))
	}

	// Falling to depth 2 closes Inner (opened at depth 2).
	got = popClosedCsScopes(got, 2)
	if len(got) != 2 || got[len(got)-1].typeName != "Outer" {
		t.Fatalf("depth 2: got %+v, want top=Outer", got)
	}

	// Falling to depth 0 closes everything.
	got = popClosedCsScopes(got, 0)
	if len(got) != 0 {
		t.Fatalf("depth 0: got %d scopes, want 0", len(got))
	}
}
