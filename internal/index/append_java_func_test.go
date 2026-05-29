package index

import "testing"

func TestAppendJavaFunc(t *testing.T) {
	st := &javaParseState{relPath: "src/main/java/p/Foo.java", pkg: "p", lineNum: 7,
		scopes: []javaScope{{typeName: "Foo"}}}
	appendJavaFunc(st, "public int add() {", "add")
	if len(st.functions) != 1 {
		t.Fatalf("got %d functions, want 1", len(st.functions))
	}
	fn := st.functions[0]
	if fn.Name != "add" || !fn.Exported || fn.StartLine != 7 || fn.File != "src/main/java/p/Foo.java" {
		t.Errorf("function = %+v", fn)
	}
	if fn.QualifiedName != "p.Foo.add" {
		t.Errorf("qualified = %q, want p.Foo.add", fn.QualifiedName)
	}
}

func TestAppendJavaFuncNonPublic(t *testing.T) {
	st := &javaParseState{relPath: "Foo.java"}
	appendJavaFunc(st, "void helper() {", "helper")
	if len(st.functions) != 1 || st.functions[0].Exported {
		t.Errorf("expected non-exported helper, got %+v", st.functions)
	}
}
