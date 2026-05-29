package index

import "testing"

func TestDispatchJavaLineType(t *testing.T) {
	st := &javaParseState{}
	dispatchJavaLine(st, "public class Foo {")
	if len(st.scopes) != 1 || st.scopes[0].typeName != "Foo" {
		t.Errorf("scopes = %+v, want one Foo type", st.scopes)
	}
}

func TestDispatchJavaLineMethod(t *testing.T) {
	st := &javaParseState{relPath: "Foo.java"}
	dispatchJavaLine(st, "public void run() {")
	if len(st.functions) != 1 || st.functions[0].Name != "run" {
		t.Errorf("functions = %+v, want one run", st.functions)
	}
}

func TestDispatchJavaLineOther(t *testing.T) {
	st := &javaParseState{}
	dispatchJavaLine(st, "int x = 1;")
	if len(st.scopes) != 0 || len(st.functions) != 0 {
		t.Error("non-declaration line should not change state")
	}
}
