package index

import "testing"

func TestAppendCsFunc(t *testing.T) {
	st := &csParseState{relPath: "src/Calculator.cs", fileNs: "Com.Example", lineNum: 7,
		scopes: []csScope{{typeName: "Calculator"}}}
	appendCsFunc(st, "public int Add(int a, int b) {", "Add")
	if len(st.functions) != 1 {
		t.Fatalf("got %d functions, want 1", len(st.functions))
	}
	fn := st.functions[0]
	if fn.Name != "Add" || !fn.Exported || fn.StartLine != 7 || fn.File != "src/Calculator.cs" {
		t.Errorf("function = %+v", fn)
	}
	if fn.QualifiedName != "Com.Example.Calculator.Add" {
		t.Errorf("qualified = %q, want Com.Example.Calculator.Add", fn.QualifiedName)
	}
}

func TestAppendCsFuncNonPublic(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	appendCsFunc(st, "void Helper() {", "Helper")
	if len(st.functions) != 1 || st.functions[0].Exported {
		t.Errorf("expected non-exported Helper, got %+v", st.functions)
	}
}
