package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCsIndexerIndexBasic(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `namespace P;

public class Foo
{
    public int A() { return 1; }
    public int B() { return 2; }
}
`
	if err := os.WriteFile(filepath.Join(srcDir, "Foo.cs"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	wantFile := filepath.Join("src", "Foo.cs")
	for _, f := range funcs {
		if f.File != wantFile {
			t.Errorf("File = %q, want %q", f.File, wantFile)
		}
	}
}

func TestCsIndexerIndexSkipsTestProject(t *testing.T) {
	dir := t.TempDir()
	mainDir := filepath.Join(dir, "App")
	testDir := filepath.Join(dir, "App.Tests")
	if err := os.MkdirAll(mainDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mainDir, "Foo.cs"),
		[]byte("namespace P;\npublic class Foo {\n    public int A() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "FooTests.cs"),
		[]byte("namespace P;\npublic class FooTests {\n    public void TestA() { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "A" {
		t.Fatalf("expected only [A], got %+v", funcs)
	}
}

func TestCsIndexerIndexSkipsObjBin(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Main.cs"),
		[]byte("public class Main {\n    public static void Run() { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	objDir := filepath.Join(dir, "obj", "Debug")
	if err := os.MkdirAll(objDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(objDir, "Gen.cs"),
		[]byte("public class Gen {\n    public void G() { }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (obj/ skipped), got %d: %+v", len(funcs), funcs)
	}
}

func TestCsIndexerIndexEmpty(t *testing.T) {
	dir := t.TempDir()
	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}
}

func TestCsIndexerImplementsIndexer(t *testing.T) {
	var _ Indexer = &CsIndexer{}
}

// TestCsIndexerIndexWalkError covers the walk-error branch (err != nil in the
// WalkFunc): an unreadable subdirectory yields a walk error that is swallowed,
// and indexing of a readable sibling file still succeeds.
func TestCsIndexerIndexWalkError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Keep.cs"),
		[]byte("public class Keep {\n    public int A() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(bad, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(bad, 0o755) })

	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	found := false
	for _, fn := range funcs {
		if fn.Name == "A" {
			found = true
		}
	}
	if !found {
		t.Error("expected method A() to be indexed despite unreadable subdir")
	}
}

// TestCsIndexerIndexRespectsIgnore covers the .tsmignore branches: a matched
// directory is skipped via SkipDir and a matched file is skipped via return nil.
func TestCsIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "Keep.cs"),
		[]byte("public class Keep {\n    public int A() { return 1; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Ignored.cs"),
		[]byte("public class Ignored {\n    public int Z() { return 9; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	skipDir := filepath.Join(dir, "vendored")
	if err := os.MkdirAll(skipDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skipDir, "Vend.cs"),
		[]byte("public class Vend {\n    public int V() { return 0; }\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"),
		[]byte("Ignored.cs\nvendored\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &CsIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("CsIndexer.Index: %v", err)
	}
	names := map[string]bool{}
	for _, fn := range funcs {
		names[fn.Name] = true
	}
	if !names["A"] {
		t.Error("expected method A() from Keep.cs to be indexed")
	}
	if names["Z"] {
		t.Error("Ignored.cs should have been skipped via .tsmignore")
	}
	if names["V"] {
		t.Error("vendored/ should have been skipped via .tsmignore")
	}
}
