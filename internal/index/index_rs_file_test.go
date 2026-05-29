package index

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRs(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestIndexRsFileFreeFunctions(t *testing.T) {
	dir := t.TempDir()
	content := `pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn helper() {
    println!("hi");
}
`
	abs := writeRs(t, dir, "lib.rs", content)
	funcs := indexRsFile("lib.rs", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}

	var add, helper bool
	for _, f := range funcs {
		switch f.Name {
		case "add":
			add = true
			if !f.Exported {
				t.Error("add should be exported (pub)")
			}
			if f.StartLine != 1 {
				t.Errorf("add StartLine = %d, want 1", f.StartLine)
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
		t.Errorf("missing functions: add=%v helper=%v", add, helper)
	}
}

func TestIndexRsFileImplMethods(t *testing.T) {
	dir := t.TempDir()
	content := `struct Counter {
    n: i32,
}

impl Counter {
    pub fn new() -> Self {
        Counter { n: 0 }
    }

    pub fn inc(&mut self) {
        self.n += 1;
    }
}
`
	abs := writeRs(t, dir, "counter.rs", content)
	funcs := indexRsFile("counter.rs", abs)

	if len(funcs) != 2 {
		t.Fatalf("got %d funcs, want 2: %+v", len(funcs), funcs)
	}
	for _, f := range funcs {
		if f.QualifiedName != "Counter."+f.Name {
			t.Errorf("%s qualified = %q, want Counter.%s", f.Name, f.QualifiedName, f.Name)
		}
	}
}

func TestIndexRsFileSkipsCfgTest(t *testing.T) {
	dir := t.TempDir()
	content := `pub fn real() -> i32 {
    42
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(real(), 42);
    }
}
`
	abs := writeRs(t, dir, "lib.rs", content)
	funcs := indexRsFile("lib.rs", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1 (cfg(test) skipped): %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "real" {
		t.Errorf("indexed %q, want real", funcs[0].Name)
	}
}

func TestIndexRsFileMissingFile(t *testing.T) {
	// os.Open fails on a nonexistent path -> returns nil functions.
	funcs := indexRsFile("missing.rs", filepath.Join(t.TempDir(), "does-not-exist.rs"))
	if funcs != nil {
		t.Errorf("expected nil for unreadable file, got %+v", funcs)
	}
}

func TestIndexRsFileNestedModule(t *testing.T) {
	dir := t.TempDir()
	content := `pub mod math {
    pub fn square(x: i32) -> i32 {
        x * x
    }
}
`
	abs := writeRs(t, dir, "lib.rs", content)
	funcs := indexRsFile("lib.rs", abs)

	if len(funcs) != 1 {
		t.Fatalf("got %d funcs, want 1: %+v", len(funcs), funcs)
	}
	if funcs[0].QualifiedName != "math.square" {
		t.Errorf("qualified = %q, want math.square", funcs[0].QualifiedName)
	}
}
