package index

import "testing"

func TestRsParseStateZeroValue(t *testing.T) {
	st := &rsParseState{relDir: "src", relPath: "src/lib.rs"}
	if st.depth != 0 || st.lineNum != 0 || st.pendingCfgTest {
		t.Errorf("unexpected zero-value state: %+v", st)
	}
	if len(st.functions) != 0 || len(st.scopes) != 0 {
		t.Error("expected empty functions and scopes")
	}
	if st.relDir != "src" || st.relPath != "src/lib.rs" {
		t.Errorf("relDir/relPath = %q/%q", st.relDir, st.relPath)
	}
}
