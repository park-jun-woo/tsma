package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJavaIndexerIndexBasic(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src", "main", "java", "p")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `package p;

public class Foo {
    public int a() { return 1; }
    public int b() { return 2; }
}
`
	if err := os.WriteFile(filepath.Join(srcDir, "Foo.java"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("JavaIndexer.Index: %v", err)
	}
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	wantFile := filepath.Join("src", "main", "java", "p", "Foo.java")
	for _, f := range funcs {
		if f.File != wantFile {
			t.Errorf("File = %q, want %q", f.File, wantFile)
		}
	}
}

func TestJavaIndexerIndexSkipsTestTree(t *testing.T) {
	dir := t.TempDir()
	mainDir := filepath.Join(dir, "src", "main", "java", "p")
	testDir := filepath.Join(dir, "src", "test", "java", "p")
	if err := os.MkdirAll(mainDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainDir, "Foo.java"),
		[]byte("package p;\npublic class Foo {\n    public int a() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "FooTest.java"),
		[]byte("package p;\npublic class FooTest {\n    public void testA() { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("JavaIndexer.Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "a" {
		t.Fatalf("expected only [a], got %+v", funcs)
	}
}

func TestJavaIndexerIndexSkipsTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Main.java"),
		[]byte("public class Main {\n    public static void main(String[] a) { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	targetDir := filepath.Join(dir, "target", "classes")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "Gen.java"),
		[]byte("public class Gen {\n    public void g() { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("JavaIndexer.Index: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (target/ skipped), got %d: %+v", len(funcs), funcs)
	}
}

func TestJavaIndexerIndexEmpty(t *testing.T) {
	dir := t.TempDir()
	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("JavaIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}
}

func TestJavaIndexerImplementsIndexer(t *testing.T) {
	var _ Indexer = &JavaIndexer{}
}

// TestJavaIndexerIndexWalkError covers the walk-error branch (lines 18-19):
// an unreadable subdirectory makes filepath.Walk invoke the callback with a
// non-nil err, which is swallowed so indexing still succeeds.
func TestJavaIndexerIndexWalkError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; unreadable-dir walk error does not apply")
	}
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "Keep.java"),
		[]byte("public class Keep {\n    public int a() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	sub := filepath.Join(dir, "locked")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "Inner.java"),
		[]byte("public class Inner {\n    public int b() { return 2; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("expected indexing to succeed despite walk error, got: %v", err)
	}
	found := false
	for _, fn := range funcs {
		if fn.Name == "a" {
			found = true
		}
	}
	if !found {
		t.Error("expected method a() from Keep.java to still be indexed")
	}
}

// TestJavaIndexerIndexRespectsIgnore covers the .tsmignore branches: a matched
// directory is skipped via SkipDir (lines 22-24) and a matched file is skipped
// via return nil (line 26).
func TestJavaIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "Keep.java"),
		[]byte("public class Keep {\n    public int a() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A file that will be ignored by name.
	if err := os.WriteFile(filepath.Join(dir, "Ignored.java"),
		[]byte("public class Ignored {\n    public int z() { return 9; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A whole directory that will be ignored.
	skipDir := filepath.Join(dir, "vendored")
	if err := os.MkdirAll(skipDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skipDir, "Vend.java"),
		[]byte("public class Vend {\n    public int v() { return 0; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"),
		[]byte("Ignored.java\nvendored\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &JavaIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("JavaIndexer.Index: %v", err)
	}
	names := map[string]bool{}
	for _, fn := range funcs {
		names[fn.Name] = true
	}
	if !names["a"] {
		t.Error("expected method a() from Keep.java to be indexed")
	}
	if names["z"] {
		t.Error("Ignored.java should have been skipped via .tsmignore")
	}
	if names["v"] {
		t.Error("vendored/ should have been skipped via .tsmignore")
	}
}
