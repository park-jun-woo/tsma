package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// fakeRsXML is a canned Rust parse tree: a #[cfg(test)] module whose body calls
// target_fn() — exercising rsCfgTestBodies for the in-file path while the same
// call_expression also feeds the whole-file walk of ingestRsTestFile.
const fakeRsXML = `<source_file srow="0" scol="0" erow="8" ecol="0">
 <attribute_item srow="0" scol="0" erow="0" ecol="12">
  <identifier srow="0" scol="2" erow="0" ecol="5">cfg</identifier>
  <identifier srow="0" scol="6" erow="0" ecol="10">test</identifier>
 </attribute_item>
 <mod_item srow="1" scol="0" erow="7" ecol="1">
  <declaration_list field="body" srow="1" scol="10" erow="7" ecol="1">
   <call_expression srow="3" scol="8" erow="3" ecol="19">
    <identifier field="function" srow="3" scol="8" erow="3" ecol="17">target_fn</identifier>
   </call_expression>
  </declaration_list>
 </mod_item>
</source_file>`

// fakeRsEmptyXML is a canned Rust parse tree with no test module and no calls,
// so nothing is ever indexed from it.
const fakeRsEmptyXML = `<source_file srow="0" scol="0" erow="1" ecol="0">
 <identifier srow="0" scol="0" erow="0" ecol="1">x</identifier>
</source_file>`

// writeRsFakeCrate lays out a minimal crate (src/lib.rs without an in-file test
// module) under a fresh root and returns it.
func writeRsFakeCrate(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	abs := filepath.Join(root, "src", "lib.rs")
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte("pub fn target_fn() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestBuildRsTestIndexFakeCLI covers the CLI-present branches without a real
// tree-sitter install: the in-file #[cfg(test)] module and a tests/*.rs
// integration file both contribute references, and a parse yielding no
// references gives nil (len==0).
func TestBuildRsTestIndexFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_RUST_GRAMMAR", fakeRsXML)
	root := writeRsFakeCrate(t)
	if err := os.MkdirAll(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "it.rs"), []byte("t\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := BuildRsTestIndex(root, "src/lib.rs")
	if idx == nil {
		t.Fatal("expected a populated index for src/lib.rs")
	}
	got := idx.refs["target_fn"]
	want := []string{"src/lib.rs", filepath.Join("tests", "it.rs")}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("refs[target_fn] = %v, want %v", got, want)
	}

	// A parse that yields no cfg(test) bodies and no tests/ refs -> nil.
	useFakeTreeSitter(t, "TSMA_RUST_GRAMMAR", fakeRsEmptyXML)
	if got := BuildRsTestIndex(writeRsFakeCrate(t), "src/lib.rs"); got != nil {
		t.Errorf("expected nil index when nothing is referenced, got %+v", got)
	}
}

// TestIngestRsInFileFakeCLI covers the successful parse branch: every name the
// in-file #[cfg(test)] module calls gets a back-reference to the source file.
func TestIngestRsInFileFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeRsXML)
	idx := &RsTestIndex{refs: map[string][]string{}}
	ingestRsInFile(idx, script, "", "/x/src/lib.rs", "src/lib.rs")
	if got := idx.refs["target_fn"]; len(got) != 1 || got[0] != "src/lib.rs" {
		t.Errorf("refs[target_fn] = %v, want [src/lib.rs]", got)
	}
}

// TestIngestRsTestFileFakeCLI covers the successful parse branch: every name
// the integration file calls gets a back-reference to that file.
func TestIngestRsTestFileFakeCLI(t *testing.T) {
	script := writeFakeTreeSitter(t, fakeRsXML)
	idx := &RsTestIndex{refs: map[string][]string{}}
	ingestRsTestFile(idx, script, "", "/x/tests/it.rs", "tests/it.rs")
	if got := idx.refs["target_fn"]; len(got) != 1 || got[0] != "tests/it.rs" {
		t.Errorf("refs[target_fn] = %v, want [tests/it.rs]", got)
	}
}

// TestRsFuncMatcherMatchFuncFakeCLI covers the content-index branches of
// MatchFunc: a referenced name attributes to its own source file via the
// in-file module, and an unreferenced name falls through to the filename
// fallback (here: not found — no in-file tests, no tests/lib.rs).
func TestRsFuncMatcherMatchFuncFakeCLI(t *testing.T) {
	useFakeTreeSitter(t, "TSMA_RUST_GRAMMAR", fakeRsXML)
	root := writeRsFakeCrate(t)
	m := &RsFuncMatcher{}

	// Content hit: target_fn is called by the in-file #[cfg(test)] module.
	fn := &model.Function{Name: "target_fn", File: "src/lib.rs"}
	tm, ok := m.MatchFunc(root, fn)
	if !ok || len(tm.Files) != 1 || tm.Files[0] != "src/lib.rs" {
		t.Errorf("content hit: ok=%v tm=%+v, want Files=[src/lib.rs]", ok, tm)
	}

	// Content miss with a populated index -> filename fallback -> not found.
	miss := &model.Function{Name: "absent", File: "src/lib.rs"}
	if _, ok := m.MatchFunc(root, miss); ok {
		t.Error("content miss without a conventional test file should report not found")
	}
}
