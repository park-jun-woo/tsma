package index

import "testing"

func TestDispatchCsLineFileScopedNamespace(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	dispatchCsLine(st, "namespace Com.Example;")
	if st.fileNs != "Com.Example" {
		t.Errorf("fileNs = %q, want Com.Example", st.fileNs)
	}
	if len(st.pending) != 0 || len(st.scopes) != 0 {
		t.Errorf("file-scoped namespace must not push/pend a scope, pending=%+v scopes=%+v", st.pending, st.scopes)
	}
}

func TestDispatchCsLineBlockNamespace(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	dispatchCsLine(st, "namespace Com.Example {")
	if len(st.pending) != 1 || st.pending[0] != "Com.Example" {
		t.Errorf("block namespace should be pending, got %+v", st.pending)
	}
}

func TestDispatchCsLineType(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	dispatchCsLine(st, "public class Calculator {")
	if len(st.pending) != 1 || st.pending[0] != "Calculator" {
		t.Errorf("class should be pending, got %+v", st.pending)
	}
}

func TestDispatchCsLineMethod(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs", scopes: []csScope{{typeName: "Calculator"}}}
	dispatchCsLine(st, "public int Add(int a, int b) {")
	if len(st.functions) != 1 || st.functions[0].Name != "Add" {
		t.Errorf("method should be recorded, got %+v", st.functions)
	}
}
