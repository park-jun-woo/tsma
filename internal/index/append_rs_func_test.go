package index

import "testing"

func TestAppendRsFunc(t *testing.T) {
	st := &rsParseState{relPath: "src/lib.rs", lineNum: 5}
	appendRsFunc(st, "pub fn add() {")
	if len(st.functions) != 1 {
		t.Fatalf("got %d functions, want 1", len(st.functions))
	}
	fn := st.functions[0]
	if fn.Name != "add" || !fn.Exported || fn.StartLine != 5 || fn.File != "src/lib.rs" {
		t.Errorf("function = %+v", fn)
	}
}

func TestAppendRsFuncSkipsCfgTest(t *testing.T) {
	st := &rsParseState{relPath: "src/lib.rs", pendingCfgTest: true}
	appendRsFunc(st, "fn it_works() {")
	if len(st.functions) != 0 {
		t.Errorf("expected cfg(test) fn to be skipped, got %d", len(st.functions))
	}
}
