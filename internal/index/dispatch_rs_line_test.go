package index

import "testing"

func TestDispatchRsLineImpl(t *testing.T) {
	st := &rsParseState{}
	dispatchRsLine(st, "impl Foo {", "impl Foo {")
	if len(st.scopes) != 1 || st.scopes[0].receiver != "Foo" {
		t.Errorf("scopes = %+v, want one Foo receiver", st.scopes)
	}
}

func TestDispatchRsLineMod(t *testing.T) {
	st := &rsParseState{}
	dispatchRsLine(st, "pub mod util {", "pub mod util {")
	if len(st.scopes) != 1 || st.scopes[0].module != "util" {
		t.Errorf("scopes = %+v, want one util module", st.scopes)
	}
}

func TestDispatchRsLineFn(t *testing.T) {
	st := &rsParseState{relPath: "lib.rs"}
	dispatchRsLine(st, "fn helper() {", "fn helper() {")
	if len(st.functions) != 1 || st.functions[0].Name != "helper" {
		t.Errorf("functions = %+v, want one helper", st.functions)
	}
}

func TestDispatchRsLineOther(t *testing.T) {
	st := &rsParseState{}
	dispatchRsLine(st, "let x = 1;", "let x = 1;")
	if len(st.scopes) != 0 || len(st.functions) != 0 {
		t.Error("non-declaration line should not change state")
	}
}
