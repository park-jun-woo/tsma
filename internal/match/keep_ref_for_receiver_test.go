package match

import "testing"

func TestKeepRefForReceiver(t *testing.T) {
	cases := []struct {
		name     string
		ref      testRef
		fnRecv   string
		multiple bool
		want     bool
	}{
		{"exact", testRef{Receiver: "GoFile"}, "GoFile", true, true},
		{"mismatch", testRef{Receiver: "CSharpFile"}, "GoFile", false, false},
		{"unknown-single", testRef{Receiver: ""}, "GoFile", false, true},
		{"unknown-multiple", testRef{Receiver: ""}, "GoFile", true, false},
	}
	for _, c := range cases {
		if got := keepRefForReceiver(c.ref, c.fnRecv, c.multiple); got != c.want {
			t.Errorf("%s: keepRefForReceiver = %v, want %v", c.name, got, c.want)
		}
	}
}
