package index

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCs(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestIndexCsFileMethods(t *testing.T) {
	dir := t.TempDir()
	content := `namespace Com.Example;

public class Calculator
{
    public int Add(int a, int b)
    {
        return a + b;
    }

    private int Helper(int x)
    {
        return x * 2;
    }
}
`
	abs := writeCs(t, dir, "Calculator.cs", content)
	funcs := indexCsFile("Calculator.cs", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
	var add, helper bool
	for _, f := range funcs {
		switch f.Name {
		case "Add":
			add = true
			if !f.Exported {
				t.Error("Add should be exported (public)")
			}
			if f.QualifiedName != "Com.Example.Calculator.Add" {
				t.Errorf("Add qualified = %q", f.QualifiedName)
			}
			if f.EndLine < f.StartLine {
				t.Errorf("Add EndLine %d < StartLine %d", f.EndLine, f.StartLine)
			}
		case "Helper":
			helper = true
			if f.Exported {
				t.Error("Helper should not be exported")
			}
		}
	}
	if !add || !helper {
		t.Errorf("missing methods: Add=%v Helper=%v", add, helper)
	}
}

func TestIndexCsFileSameLineBody(t *testing.T) {
	dir := t.TempDir()
	content := `namespace P;
public class Foo {
    public int A() { return 1; }
    public int B() { return 2; }
}
`
	abs := writeCs(t, dir, "Foo.cs", content)
	funcs := indexCsFile("Foo.cs", abs)
	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
}

func TestIndexCsFileSkipsAttributes(t *testing.T) {
	dir := t.TempDir()
	content := `namespace P;

public class Svc
{
    [Obsolete]
    public string Describe()
    {
        return "svc";
    }
}
`
	abs := writeCs(t, dir, "Svc.cs", content)
	funcs := indexCsFile("Svc.cs", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1: %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "Describe" {
		t.Errorf("indexed %q, want Describe", funcs[0].Name)
	}
}

func TestIndexCsFileBlockNamespaceNested(t *testing.T) {
	dir := t.TempDir()
	content := `namespace Com.Example
{
    public class Outer
    {
        public class Inner
        {
            public void Run()
            {
                DoWork();
            }
        }

        public void OuterMethod()
        {
        }
    }
}
`
	abs := writeCs(t, dir, "Outer.cs", content)
	funcs := indexCsFile("Outer.cs", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
	var run, outer bool
	for _, f := range funcs {
		switch f.Name {
		case "Run":
			run = true
			if f.QualifiedName != "Com.Example.Outer.Inner.Run" {
				t.Errorf("Run qualified = %q, want Com.Example.Outer.Inner.Run", f.QualifiedName)
			}
		case "OuterMethod":
			outer = true
			if f.QualifiedName != "Com.Example.Outer.OuterMethod" {
				t.Errorf("OuterMethod qualified = %q, want Com.Example.Outer.OuterMethod", f.QualifiedName)
			}
		}
	}
	if !run || !outer {
		t.Errorf("missing methods: Run=%v OuterMethod=%v", run, outer)
	}
}

func TestIndexCsFileConstructorAndControl(t *testing.T) {
	dir := t.TempDir()
	content := `namespace P;

public class Widget
{
    public Widget(int n)
    {
        if (n > 0)
        {
            Init();
        }
        foreach (var x in items)
        {
        }
    }
}
`
	abs := writeCs(t, dir, "Widget.cs", content)
	funcs := indexCsFile("Widget.cs", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1 (constructor only, control skipped): %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "Widget" {
		t.Errorf("indexed %q, want Widget", funcs[0].Name)
	}
}

func TestIndexCsFileMissingFile(t *testing.T) {
	funcs := indexCsFile("missing.cs", filepath.Join(t.TempDir(), "does-not-exist.cs"))
	if funcs != nil {
		t.Errorf("expected nil for unreadable file, got %+v", funcs)
	}
}
