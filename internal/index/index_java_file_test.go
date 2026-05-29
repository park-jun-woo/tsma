package index

import (
	"os"
	"path/filepath"
	"testing"
)

func writeJava(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestIndexJavaFileMethods(t *testing.T) {
	dir := t.TempDir()
	content := `package com.example;

public class Calculator {

    public int add(int a, int b) {
        return a + b;
    }

    private int helper(int x) {
        return x * 2;
    }
}
`
	abs := writeJava(t, dir, "Calculator.java", content)
	funcs := indexJavaFile("Calculator.java", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
	var add, helper bool
	for _, f := range funcs {
		switch f.Name {
		case "add":
			add = true
			if !f.Exported {
				t.Error("add should be exported (public)")
			}
			if f.QualifiedName != "com.example.Calculator.add" {
				t.Errorf("add qualified = %q", f.QualifiedName)
			}
			if f.EndLine < f.StartLine {
				t.Errorf("add EndLine %d < StartLine %d", f.EndLine, f.StartLine)
			}
		case "helper":
			helper = true
			if f.Exported {
				t.Error("helper should not be exported")
			}
		}
	}
	if !add || !helper {
		t.Errorf("missing methods: add=%v helper=%v", add, helper)
	}
}

func TestIndexJavaFileSkipsAnnotations(t *testing.T) {
	dir := t.TempDir()
	content := `package p;

public class Svc {
    @Override
    public String toString() {
        return "svc";
    }
}
`
	abs := writeJava(t, dir, "Svc.java", content)
	funcs := indexJavaFile("Svc.java", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1: %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "toString" {
		t.Errorf("indexed %q, want toString", funcs[0].Name)
	}
}

func TestIndexJavaFileNestedClass(t *testing.T) {
	dir := t.TempDir()
	content := `package p;

public class Outer {
    static class Inner {
        void run() {
            doWork();
        }
    }

    void outerMethod() {
    }
}
`
	abs := writeJava(t, dir, "Outer.java", content)
	funcs := indexJavaFile("Outer.java", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
	var run, outer bool
	for _, f := range funcs {
		switch f.Name {
		case "run":
			run = true
			if f.QualifiedName != "p.Outer.Inner.run" {
				t.Errorf("run qualified = %q, want p.Outer.Inner.run", f.QualifiedName)
			}
		case "outerMethod":
			outer = true
			if f.QualifiedName != "p.Outer.outerMethod" {
				t.Errorf("outerMethod qualified = %q, want p.Outer.outerMethod", f.QualifiedName)
			}
		}
	}
	if !run || !outer {
		t.Errorf("missing methods: run=%v outerMethod=%v", run, outer)
	}
}

func TestIndexJavaFileConstructorAndControl(t *testing.T) {
	dir := t.TempDir()
	content := `package p;

public class Widget {
    public Widget(int n) {
        if (n > 0) {
            init();
        }
        for (int i = 0; i < n; i++) {
        }
    }
}
`
	abs := writeJava(t, dir, "Widget.java", content)
	funcs := indexJavaFile("Widget.java", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1 (constructor only, control skipped): %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "Widget" {
		t.Errorf("indexed %q, want Widget", funcs[0].Name)
	}
}

func TestIndexJavaFileMissingFile(t *testing.T) {
	funcs := indexJavaFile("missing.java", filepath.Join(t.TempDir(), "does-not-exist.java"))
	if funcs != nil {
		t.Errorf("expected nil for unreadable file, got %+v", funcs)
	}
}
