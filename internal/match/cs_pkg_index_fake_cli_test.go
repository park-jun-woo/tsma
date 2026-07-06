package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeCsXML is a canned C# parse tree containing one bare invocation
// (TargetFn()) and one object creation (new TargetType()) — the two node kinds
// collectCsCalledNames harvests.
const fakeCsXML = `<compilation_unit srow="0" scol="0" erow="4" ecol="0">
 <invocation_expression srow="1" scol="0" erow="1" ecol="10">
  <identifier field="function" srow="1" scol="0" erow="1" ecol="8">TargetFn</identifier>
 </invocation_expression>
 <object_creation_expression srow="2" scol="0" erow="2" ecol="20">
  <identifier field="type" srow="2" scol="4" erow="2" ecol="14">TargetType</identifier>
 </object_creation_expression>
</compilation_unit>`

// writeCsFakeProject lays out a minimal C# project (Proj/Foo.cs with a parallel
// Proj.Tests/FooTests.cs) under a fresh root and returns it.
func writeCsFakeProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range map[string]string{
		"Proj/Foo.cs":            "public class Foo {}\n",
		"Proj.Tests/FooTests.cs": "public class FooTests {}\n",
	} {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// TestBuildCsPkgTestIndexFakeCLI covers the CLI-present branches without a real
// tree-sitter install: the parallel *.Tests dir yields a populated index, and a
// source dir with no test candidates yields nil (len==0).
func TestBuildCsPkgTestIndexFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_CSHARP_GRAMMAR", fakeCsXML)
	root := writeCsFakeProject(t)

	idx := BuildCsPkgTestIndex(root, "Proj")
	if idx == nil {
		t.Fatal("expected a populated index for Proj")
	}
	want := filepath.Join("Proj.Tests", "FooTests.cs")
	if got := idx.refs["TargetFn"]; len(got) != 1 || got[0] != want {
		t.Errorf("refs[TargetFn] = %v, want [%s]", got, want)
	}
	if got := idx.refs["TargetType"]; len(got) != 1 || got[0] != want {
		t.Errorf("refs[TargetType] = %v, want [%s]", got, want)
	}

	// No test dir and no in-dir test files -> nothing ingested -> nil.
	if got := BuildCsPkgTestIndex(root, "NoSuch"); got != nil {
		t.Errorf("expected nil index for a dir without test candidates, got %+v", got)
	}
}

// TestIngestCsTestDirFakeCLI covers the per-entry branches deterministically: a
// subdirectory (IsDir skip), a non-test .cs file (isCsTestFile skip), and a
// *Tests.cs file (ingested).
func TestIngestCsTestDirFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeCsXML)
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Plain.cs", "FooTests.cs"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestDir(idx, dir, ".", script, "")
	if got := idx.refs["TargetFn"]; len(got) != 1 || got[0] != "FooTests.cs" {
		t.Errorf("refs[TargetFn] = %v, want [FooTests.cs] (only the test file ingested)", got)
	}
}

// TestIngestCsTestFileFakeCLI covers the successful parse + ref recording
// branch without a real CLI.
func TestIngestCsTestFileFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeCsXML)
	idx := &CsPkgTestIndex{refs: map[string][]string{}}
	ingestCsTestFile(idx, script, "", "/x/FooTests.cs", "FooTests.cs")
	if got := idx.refs["TargetFn"]; len(got) != 1 || got[0] != "FooTests.cs" {
		t.Errorf("refs[TargetFn] = %v, want [FooTests.cs]", got)
	}
}

// TestCsFuncMatcherMatchFuncFakeCLI covers the content-index branches of
// MatchFunc: a referenced name attributes to the calling test file, and an
// unreferenced name falls through to the filename fallback (here: not found).
func TestCsFuncMatcherMatchFuncFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_CSHARP_GRAMMAR", fakeCsXML)
	root := writeCsFakeProject(t)
	m := &CsFuncMatcher{}

	// Content hit: TargetFn is called by Proj.Tests/FooTests.cs.
	fn := &model.Function{Name: "TargetFn", File: "Proj/Foo.cs"}
	tm, ok := m.MatchFunc(root, fn)
	want := filepath.Join("Proj.Tests", "FooTests.cs")
	if !ok || len(tm.Files) != 1 || tm.Files[0] != want {
		t.Errorf("content hit: ok=%v tm=%+v, want Files=[%s]", ok, tm, want)
	}

	// Content miss with a populated index -> filename fallback (no
	// conventional BareTests.cs exists) -> not found.
	miss := &model.Function{Name: "Absent", File: "Proj/Bare.cs"}
	if _, ok := m.MatchFunc(root, miss); ok {
		t.Error("content miss without a conventional test file should report not found")
	}
}
