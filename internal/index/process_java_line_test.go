package index

import "testing"

func TestProcessJavaLineMethod(t *testing.T) {
	st := &javaParseState{relPath: "Foo.java"}
	processJavaLine(st, "package p;")
	processJavaLine(st, "public class Foo {")
	processJavaLine(st, "    public int add() {")
	processJavaLine(st, "        return 1;")
	processJavaLine(st, "    }")
	processJavaLine(st, "}")

	if st.pkg != "p" {
		t.Errorf("pkg = %q, want p", st.pkg)
	}
	if len(st.functions) != 1 || st.functions[0].Name != "add" {
		t.Fatalf("functions = %+v, want one add", st.functions)
	}
	if st.functions[0].QualifiedName != "p.Foo.add" {
		t.Errorf("qualified = %q, want p.Foo.add", st.functions[0].QualifiedName)
	}
	if st.depth != 0 {
		t.Errorf("depth = %d, want 0 after closing braces", st.depth)
	}
}

func TestProcessJavaLineSkipsAnnotation(t *testing.T) {
	st := &javaParseState{relPath: "Foo.java"}
	processJavaLine(st, "@Override")
	processJavaLine(st, "public String name() {")
	processJavaLine(st, "}")

	if len(st.functions) != 1 || st.functions[0].Name != "name" {
		t.Fatalf("functions = %+v, want one name", st.functions)
	}
}
