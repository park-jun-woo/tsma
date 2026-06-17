package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// writeCsTestFixture writes a minimal xUnit-style file calling foo() and
// constructing new Bar() into dir, returning its absolute path.
func writeCsTestFixture(t *testing.T, dir, name string) string {
	t.Helper()
	abs := filepath.Join(dir, name)
	src := "namespace P;\n" +
		"public class FooTests\n" +
		"{\n" +
		"    public void T()\n" +
		"    {\n" +
		"        Foo();\n" +
		"        var b = new Bar();\n" +
		"    }\n" +
		"}\n"
	if err := os.WriteFile(abs, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return abs
}

// TestBuildCsPkgTestIndexNoTreeSitter forces the CLI absent so the command==""
// guard returns nil (signal to fall back to filename matching).
func TestBuildCsPkgTestIndexNoTreeSitter(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	if idx := BuildCsPkgTestIndex("../../testdata/csharp", "Calc"); idx != nil {
		t.Errorf("expected nil index when tree-sitter absent, got %+v", idx)
	}
}

// TestBuildCsPkgTestIndexSuccessAndEmpty covers the CLI-present branches: the
// Calc package (parallel Calc.Tests project) yields a non-nil index keyed by
// called names, and a package whose test mirror is missing yields nil (len==0).
func TestBuildCsPkgTestIndexSuccessAndEmpty(t *testing.T) {
	if !locateCsTS(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}

	idx := BuildCsPkgTestIndex("../../testdata/csharp", "Calc")
	if idx == nil {
		t.Fatal("expected non-nil index for the Calc package")
	}
	if _, ok := idx.refs["Classify"]; !ok {
		t.Errorf("index missing ref for Classify: %+v", idx.refs)
	}

	// A source package whose test mirror dir does not exist -> empty -> nil.
	if got := BuildCsPkgTestIndex("../../testdata/csharp", "Missing"); got != nil {
		t.Errorf("expected nil for missing test mirror, got %+v", got)
	}
}

// TestIngestCsTestDirMissing covers the unreadable-directory no-op branch.
func TestIngestCsTestDirMissing(t *testing.T) {
	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestDir(idx, "/", "definitely/not/a/real/dir", "tree-sitter", "")
	if len(idx.refs) != 0 {
		t.Errorf("missing dir should leave index empty, got %+v", idx.refs)
	}
}

// TestIngestCsTestDirScan covers the three per-entry branches with a real CLI: a
// subdirectory (IsDir skip), a non-test .cs file (isCsTestFile skip), and a real
// *Tests.cs file (ingested).
func TestIngestCsTestDirScan(t *testing.T) {
	if !locateCsTS(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Helper.cs"), []byte("public class Helper {}\n"), 0o644); err != nil {
		t.Fatalf("write Helper: %v", err)
	}
	writeCsTestFixture(t, dir, "FooTests.cs")

	command := treesitter.ResolveCommand()
	grammar := treesitter.ResolveGrammar("csharp")
	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestDir(idx, dir, ".", command, grammar)

	if _, ok := idx.refs["Foo"]; !ok {
		t.Errorf("expected Foo in refs after scan: %+v", idx.refs)
	}
	if _, ok := idx.refs["Bar"]; !ok {
		t.Errorf("expected Bar in refs after scan: %+v", idx.refs)
	}
}

// TestIngestCsTestFileParseError covers the parse-failure skip branch (bad
// command -> ParseFile errors -> no-op).
func TestIngestCsTestFileParseError(t *testing.T) {
	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestFile(idx, "/nonexistent/abs/tree-sitter", "", "/nope.cs", "nope.cs")
	if len(idx.refs) != 0 {
		t.Errorf("parse error should leave index empty, got %+v", idx.refs)
	}
}

// TestIngestCsTestFileSuccess covers the successful parse + ref recording.
func TestIngestCsTestFileSuccess(t *testing.T) {
	if !locateCsTS(t) {
		t.Skip("tree-sitter CLI + c-sharp grammar not available")
	}
	dir := t.TempDir()
	abs := writeCsTestFixture(t, dir, "FooTests.cs")
	command := treesitter.ResolveCommand()
	grammar := treesitter.ResolveGrammar("csharp")

	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestFile(idx, command, grammar, abs, "FooTests.cs")
	if got := idx.refs["Foo"]; len(got) != 1 || got[0] != "FooTests.cs" {
		t.Errorf("refs[Foo] = %v, want [FooTests.cs]", got)
	}
}

// TestCsFilenameFallback covers both the found (conventional FooTests.cs exists
// in the parallel *.Tests project) and not-found branches against the testdata
// layout — pure os.Stat filename matching, no CLI needed.
func TestCsFilenameFallback(t *testing.T) {
	root := "../../testdata/csharp"

	// Found: Calc/Calculator.cs -> Calc.Tests/CalculatorTests.cs.
	fn := &model.Function{Name: "Classify", File: "Calc/Calculator.cs"}
	tm, ok := csFilenameFallback(root, fn)
	if !ok {
		t.Fatalf("expected fallback match for %s", fn.File)
	}
	if len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "Calc.Tests/CalculatorTests.cs" {
		t.Errorf("Files = %v", tm.Files)
	}
	if tm.TestFuncs != nil {
		t.Errorf("TestFuncs = %v, want nil (run whole class)", tm.TestFuncs)
	}

	// Not found: no test mirror file for this source.
	missing := &model.Function{Name: "x", File: "Calc/Nonexistent.cs"}
	if _, ok := csFilenameFallback(root, missing); ok {
		t.Error("expected no fallback for a source without a test mirror")
	}
}

// TestCsFuncMatcherMatchBranches covers MatchFunc's nil-guard and the
// idx==nil -> filename-fallback path (both found and not-found) deterministically
// by forcing tree-sitter absent, complementing the CLI-gated content tests.
func TestCsFuncMatcherMatchBranches(t *testing.T) {
	m := &CsFuncMatcher{}

	// nil function.
	if _, ok := m.MatchFunc("../../testdata/csharp", nil); ok {
		t.Error("nil function should report not found")
	}

	// Force idx==nil (no tree-sitter) -> filename fallback found.
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	fn := &model.Function{Name: "Classify", File: "Calc/Calculator.cs"}
	tm, ok := m.MatchFunc("../../testdata/csharp", fn)
	if !ok || len(tm.Files) != 1 || filepath.ToSlash(tm.Files[0]) != "Calc.Tests/CalculatorTests.cs" {
		t.Errorf("fallback path: ok=%v tm=%+v", ok, tm)
	}

	// idx==nil and no test mirror -> not found.
	missing := &model.Function{Name: "x", File: "Calc/Nonexistent.cs"}
	if _, ok := m.MatchFunc("../../testdata/csharp", missing); ok {
		t.Error("no mirror should report not found")
	}
}
