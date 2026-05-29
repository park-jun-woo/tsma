package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRsIndexerIndexBasic(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn helper() {}
`
	if err := os.WriteFile(filepath.Join(srcDir, "lib.rs"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &RsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("RsIndexer.Index: %v", err)
	}
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	for _, f := range funcs {
		if f.File != filepath.Join("src", "lib.rs") {
			t.Errorf("File = %q, want relative src/lib.rs", f.File)
		}
	}
}

func TestRsIndexerIndexSkipsTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.rs"), []byte("fn main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	targetDir := filepath.Join(dir, "target", "debug")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "gen.rs"), []byte("fn generated() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &RsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("RsIndexer.Index: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (target/ skipped), got %d", len(funcs))
	}
}

func TestRsIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.rs"), []byte("fn keep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Non-Rust file -> isRsSource false -> skipped (line 31 branch).
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Ignored FILE -> file-ignore branch returns nil.
	if err := os.WriteFile(filepath.Join(dir, "ignoreme.rs"), []byte("fn skipped() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Ignored DIR -> SkipDir branch.
	ignDir := filepath.Join(dir, "vendored")
	if err := os.MkdirAll(ignDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ignDir, "dep.rs"), []byte("fn dep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("ignoreme.rs\nvendored/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &RsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("RsIndexer.Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "keep" {
		t.Fatalf("expected only keep, got %+v", funcs)
	}
}

func TestRsIndexerIndexWalkError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; unreadable-dir walk error does not apply")
	}
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.rs"), []byte("fn keep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "locked")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "inner.rs"), []byte("fn inner() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	idx := &RsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("expected success despite walk error, got: %v", err)
	}
	found := false
	for _, fn := range funcs {
		if fn.Name == "keep" {
			found = true
		}
	}
	if !found {
		t.Error("expected keep to still be indexed")
	}
}

func TestRsIndexerIndexEmpty(t *testing.T) {
	dir := t.TempDir()
	idx := &RsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("RsIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}
}
