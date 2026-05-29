package index

import "testing"

func TestJavaParseStateZeroValue(t *testing.T) {
	st := &javaParseState{relDir: "src/main/java/p", relPath: "src/main/java/p/Foo.java"}
	if st.depth != 0 || st.lineNum != 0 || st.pkg != "" {
		t.Errorf("unexpected zero-value state: %+v", st)
	}
	if len(st.functions) != 0 || len(st.scopes) != 0 {
		t.Error("expected empty functions and scopes")
	}
	if st.relDir != "src/main/java/p" || st.relPath != "src/main/java/p/Foo.java" {
		t.Errorf("relDir/relPath = %q/%q", st.relDir, st.relPath)
	}
}
