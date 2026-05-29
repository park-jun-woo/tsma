package index

import "testing"

func TestHandleRsAttribute(t *testing.T) {
	st := &rsParseState{}
	if !handleRsAttribute(st, "#[cfg(test)]") {
		t.Error("expected true for attribute line")
	}
	if !st.pendingCfgTest {
		t.Error("expected pendingCfgTest set for #[cfg(test)]")
	}

	st2 := &rsParseState{}
	if !handleRsAttribute(st2, "#[derive(Debug)]") {
		t.Error("expected true for attribute line")
	}
	if st2.pendingCfgTest {
		t.Error("non-cfg(test) attribute should not set pendingCfgTest")
	}

	st3 := &rsParseState{}
	if handleRsAttribute(st3, "pub fn f() {") {
		t.Error("expected false for non-attribute line")
	}
}
