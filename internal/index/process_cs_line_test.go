package index

import "testing"

func TestProcessCsLineMethod(t *testing.T) {
	st := &csParseState{relPath: "Calculator.cs"}
	processCsLine(st, "namespace Com.Example;")
	processCsLine(st, "public class Calculator {")
	processCsLine(st, "    public int Add(int a, int b) {")
	processCsLine(st, "        return a + b;")
	processCsLine(st, "    }")
	processCsLine(st, "}")

	if st.fileNs != "Com.Example" {
		t.Errorf("fileNs = %q, want Com.Example", st.fileNs)
	}
	if len(st.functions) != 1 || st.functions[0].Name != "Add" {
		t.Fatalf("functions = %+v, want one Add", st.functions)
	}
	if st.functions[0].QualifiedName != "Com.Example.Calculator.Add" {
		t.Errorf("qualified = %q, want Com.Example.Calculator.Add", st.functions[0].QualifiedName)
	}
	if st.depth != 0 {
		t.Errorf("depth = %d, want 0 after closing braces", st.depth)
	}
}

func TestProcessCsLineSkipsAttribute(t *testing.T) {
	st := &csParseState{relPath: "Foo.cs"}
	processCsLine(st, "[Obsolete]")
	processCsLine(st, "public string Name() {")
	processCsLine(st, "}")

	if len(st.functions) != 1 || st.functions[0].Name != "Name" {
		t.Fatalf("functions = %+v, want one Name", st.functions)
	}
}
