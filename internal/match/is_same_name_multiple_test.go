package match

import "testing"

func TestIsSameNameMultiple(t *testing.T) {
	r := &PkgSourceReceivers{byName: map[string]map[string]struct{}{
		"Single":   {"GoFile": {}},
		"Multiple": {"GoFile": {}, "CSharpFile": {}},
	}}
	if r.isSameNameMultiple("Single") {
		t.Errorf("Single should be single (false)")
	}
	if !r.isSameNameMultiple("Multiple") {
		t.Errorf("Multiple should be multiple (true)")
	}
	// Absent name -> treated as single (false), conservative.
	if r.isSameNameMultiple("Absent") {
		t.Errorf("absent name should be single (false)")
	}
}
