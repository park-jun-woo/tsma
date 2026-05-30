package match

import "testing"

func TestAddNameReceiver(t *testing.T) {
	r := &PkgSourceReceivers{byName: make(map[string]map[string]struct{})}
	addNameReceiver(r, "M", "GoFile")
	addNameReceiver(r, "M", "GoFile") // dup: no-op
	if len(r.byName["M"]) != 1 {
		t.Fatalf("after dup add, set size = %d, want 1", len(r.byName["M"]))
	}
	addNameReceiver(r, "M", "CSharpFile")
	if len(r.byName["M"]) != 2 {
		t.Fatalf("after second distinct receiver, set size = %d, want 2", len(r.byName["M"]))
	}
	if !r.isSameNameMultiple("M") {
		t.Errorf("M should now be same-name-multiple")
	}
}
