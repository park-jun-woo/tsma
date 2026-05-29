package index

import "testing"

func TestCsParseStateZeroValue(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	if st.depth != 0 {
		t.Errorf("depth = %d, want 0", st.depth)
	}
	if st.fileNs != "" {
		t.Errorf("fileNs = %q, want empty", st.fileNs)
	}
	if len(st.scopes) != 0 {
		t.Errorf("scopes = %d, want 0", len(st.scopes))
	}
	if st.relPath != "Foo.cs" {
		t.Errorf("relPath = %q, want Foo.cs", st.relPath)
	}
}
