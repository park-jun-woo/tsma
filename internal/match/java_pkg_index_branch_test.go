package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// writeJavaTestFixture writes a minimal JUnit-style file calling foo() and
// constructing new Bar() into dir, returning its absolute path.
func writeJavaTestFixture(t *testing.T, dir, name string) string {
	t.Helper()
	abs := filepath.Join(dir, name)
	src := "package p;\n" +
		"class FooTest {\n" +
		"  void t() {\n" +
		"    foo();\n" +
		"    new Bar();\n" +
		"  }\n" +
		"}\n"
	if err := os.WriteFile(abs, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return abs
}

// TestBuildJavaPkgTestIndexNoTreeSitter forces the tree-sitter CLI absent so the
// command=="" guard returns nil (the signal to fall back to filename matching).
func TestBuildJavaPkgTestIndexNoTreeSitter(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	if idx := BuildJavaPkgTestIndex("../../testdata/java", "src/main/java/com/example/calc"); idx != nil {
		t.Errorf("expected nil index when tree-sitter absent, got %+v", idx)
	}
}

// TestBuildJavaPkgTestIndexSuccessAndEmpty covers the CLI-present branches:
// a real package mirror yields a non-nil index keyed by called names, and a
// package whose test mirror directory is missing yields nil (len==0).
func TestBuildJavaPkgTestIndexSuccessAndEmpty(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}

	idx := BuildJavaPkgTestIndex("../../testdata/java", "src/main/java/com/example/calc")
	if idx == nil {
		t.Fatal("expected non-nil index for the calc package")
	}
	if _, ok := idx.refs["add"]; !ok {
		t.Errorf("index missing ref for add: %+v", idx.refs)
	}

	// A source package whose test mirror dir does not exist -> empty -> nil.
	if got := BuildJavaPkgTestIndex("../../testdata/java", "src/main/java/com/example/missingpkg"); got != nil {
		t.Errorf("expected nil for missing test mirror, got %+v", got)
	}
}

// TestIngestJavaTestDirMissing covers the unreadable-directory no-op branch.
func TestIngestJavaTestDirMissing(t *testing.T) {
	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestDir(idx, "/", "definitely/not/a/real/dir", "tree-sitter", "")
	if len(idx.refs) != 0 {
		t.Errorf("missing dir should leave index empty, got %+v", idx.refs)
	}
}

// TestIngestJavaTestDirScan covers the three per-entry branches with a real CLI:
// a subdirectory (IsDir skip), a non-test .java file (isJavaTestFile skip), and
// a real *Test.java file (ingested).
func TestIngestJavaTestDirScan(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Helper.java"), []byte("class Helper {}\n"), 0o644); err != nil {
		t.Fatalf("write Helper: %v", err)
	}
	writeJavaTestFixture(t, dir, "FooTest.java")

	command := treesitter.ResolveCommand()
	grammar := treesitter.ResolveGrammar("java")
	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestDir(idx, dir, ".", command, grammar)

	if _, ok := idx.refs["foo"]; !ok {
		t.Errorf("expected foo in refs after scan: %+v", idx.refs)
	}
	if _, ok := idx.refs["Bar"]; !ok {
		t.Errorf("expected Bar in refs after scan: %+v", idx.refs)
	}
}

// TestIngestJavaTestFileParseError covers the parse-failure skip branch (bad
// command -> ParseFile errors -> no-op).
func TestIngestJavaTestFileParseError(t *testing.T) {
	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestFile(idx, "/nonexistent/abs/tree-sitter", "", "/nope.java", "nope.java")
	if len(idx.refs) != 0 {
		t.Errorf("parse error should leave index empty, got %+v", idx.refs)
	}
}

// TestIngestJavaTestFileSuccess covers the successful parse + ref recording.
func TestIngestJavaTestFileSuccess(t *testing.T) {
	if !locateJavaTS(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	dir := t.TempDir()
	abs := writeJavaTestFixture(t, dir, "FooTest.java")
	command := treesitter.ResolveCommand()
	grammar := treesitter.ResolveGrammar("java")

	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestFile(idx, command, grammar, abs, "FooTest.java")
	if got := idx.refs["foo"]; len(got) != 1 || got[0] != "FooTest.java" {
		t.Errorf("refs[foo] = %v, want [FooTest.java]", got)
	}
}
