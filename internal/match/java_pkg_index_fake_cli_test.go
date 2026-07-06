package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeJavaXML is a canned Java parse tree containing one method invocation
// (targetFn()) and one object creation (new TargetType()) — the two node kinds
// collectJavaCalledNames harvests.
const fakeJavaXML = `<program srow="0" scol="0" erow="4" ecol="0">
 <method_invocation srow="1" scol="0" erow="1" ecol="10">
  <identifier field="name" srow="1" scol="0" erow="1" ecol="8">targetFn</identifier>
 </method_invocation>
 <object_creation_expression srow="2" scol="0" erow="2" ecol="20">
  <type_identifier field="type" srow="2" scol="4" erow="2" ecol="14">TargetType</type_identifier>
 </object_creation_expression>
</program>`

// writeJavaFakeProject lays out a minimal Maven-style project (src/main mirror
// with a src/test FooTest.java) under a fresh root and returns it.
func writeJavaFakeProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range map[string]string{
		"src/main/java/com/x/Foo.java":     "class Foo {}\n",
		"src/test/java/com/x/FooTest.java": "class FooTest {}\n",
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

// TestBuildJavaPkgTestIndexFakeCLI covers the CLI-present branches without a
// real tree-sitter install: the JUnit mirror dir yields a populated index, and
// a package whose mirror is missing yields nil (len==0).
func TestBuildJavaPkgTestIndexFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_JAVA_GRAMMAR", fakeJavaXML)
	root := writeJavaFakeProject(t)

	idx := BuildJavaPkgTestIndex(root, "src/main/java/com/x")
	if idx == nil {
		t.Fatal("expected a populated index for src/main/java/com/x")
	}
	want := filepath.Join("src", "test", "java", "com", "x", "FooTest.java")
	if got := idx.refs["targetFn"]; len(got) != 1 || got[0] != want {
		t.Errorf("refs[targetFn] = %v, want [%s]", got, want)
	}
	if got := idx.refs["TargetType"]; len(got) != 1 || got[0] != want {
		t.Errorf("refs[TargetType] = %v, want [%s]", got, want)
	}

	// Missing test mirror -> nothing ingested -> nil.
	if got := BuildJavaPkgTestIndex(root, "src/main/java/com/none"); got != nil {
		t.Errorf("expected nil index for a missing test mirror, got %+v", got)
	}
}

// TestIngestJavaTestDirFakeCLI covers the per-entry branches deterministically:
// a subdirectory (IsDir skip), a non-test .java file (isJavaTestFile skip), and
// a *Test.java file (ingested).
func TestIngestJavaTestDirFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeJavaXML)
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Plain.java", "FooTest.java"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestDir(idx, dir, ".", script, "")
	if got := idx.refs["targetFn"]; len(got) != 1 || got[0] != "FooTest.java" {
		t.Errorf("refs[targetFn] = %v, want [FooTest.java] (only the test file ingested)", got)
	}
}

// TestIngestJavaTestFileFakeCLI covers the successful parse + ref recording
// branch without a real CLI.
func TestIngestJavaTestFileFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeJavaXML)
	idx := &JavaPkgTestIndex{refs: map[string][]string{}}
	ingestJavaTestFile(idx, script, "", "/x/FooTest.java", "FooTest.java")
	if got := idx.refs["targetFn"]; len(got) != 1 || got[0] != "FooTest.java" {
		t.Errorf("refs[targetFn] = %v, want [FooTest.java]", got)
	}
}

// TestJavaFuncMatcherMatchFuncFakeCLI covers the content-index branches of
// MatchFunc: a referenced name attributes to the calling test file, and an
// unreferenced name falls through to the filename fallback (here: not found).
func TestJavaFuncMatcherMatchFuncFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_JAVA_GRAMMAR", fakeJavaXML)
	root := writeJavaFakeProject(t)
	m := &JavaFuncMatcher{}

	// Content hit: targetFn is called by the mirror's FooTest.java.
	fn := &model.Function{Name: "targetFn", File: "src/main/java/com/x/Foo.java"}
	tm, ok := m.MatchFunc(root, fn)
	want := filepath.Join("src", "test", "java", "com", "x", "FooTest.java")
	if !ok || len(tm.Files) != 1 || tm.Files[0] != want {
		t.Errorf("content hit: ok=%v tm=%+v, want Files=[%s]", ok, tm, want)
	}

	// Content miss with a populated index -> filename fallback (no
	// conventional BareTest.java exists) -> not found.
	miss := &model.Function{Name: "absent", File: "src/main/java/com/x/Bare.java"}
	if _, ok := m.MatchFunc(root, miss); ok {
		t.Error("content miss without a conventional test file should report not found")
	}
}
