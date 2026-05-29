package match

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writePkg writes a set of named files into a fresh temp package dir and
// returns the project root and the package dir relative path.
func writePkg(t *testing.T, files map[string]string) (root, pkgDir string) {
	t.Helper()
	root = t.TempDir()
	pkgDir = filepath.Join("internal", "gen")
	abs := filepath.Join(root, pkgDir)
	if err := os.MkdirAll(abs, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(abs, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, pkgDir
}

// lookupFuncs returns the sorted TestFunc names attributed to a bare func name.
func lookupFuncs(idx *PkgTestIndex, name string) []string {
	refs, ok := idx.Lookup(name)
	if !ok {
		return nil
	}
	var out []string
	for _, r := range refs {
		out = append(out, r.TestFunc)
	}
	sort.Strings(out)
	return out
}

// TestBuildPkgTestIndexGenerateBytesSeparation verifies that a call to
// GenerateBytes is attributed only to GenerateBytes and never bleeds into
// Generate (the decisive content-aware case from the plan).
func TestBuildPkgTestIndexGenerateBytesSeparation(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"generate_write_file_test.go": `package gen
import "testing"
func TestWriteFile(t *testing.T) { Generate("a", "b") }
`,
		"generate_escrow_test.go": `package gen
import "testing"
func TestEscrow(t *testing.T) { GenerateBytes("a", nil) }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	gen := lookupFuncs(idx, "Generate")
	if len(gen) != 1 || gen[0] != "TestWriteFile" {
		t.Fatalf("Generate refs = %v, want [TestWriteFile]", gen)
	}
	genBytes := lookupFuncs(idx, "GenerateBytes")
	if len(genBytes) != 1 || genBytes[0] != "TestEscrow" {
		t.Fatalf("GenerateBytes refs = %v, want [TestEscrow]", genBytes)
	}

	// Path stored must be project-root relative, matching matcher convention.
	refs, _ := idx.Lookup("Generate")
	wantRel := filepath.Join(pkgDir, "generate_write_file_test.go")
	if refs[0].File != wantRel {
		t.Errorf("File = %q, want %q", refs[0].File, wantRel)
	}
}

// TestBuildPkgTestIndexOneToMany verifies a single source func attributed to
// multiple tests spread across multiple test files (1:N).
func TestBuildPkgTestIndexOneToMany(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"parse_happy_test.go": `package gen
import "testing"
func TestParseHappy(t *testing.T) { Parse("x") }
`,
		"parse_error_test.go": `package gen
import "testing"
func TestParseError(t *testing.T) { Parse("") }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	got := lookupFuncs(idx, "Parse")
	want := []string{"TestParseError", "TestParseHappy"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Parse refs = %v, want %v", got, want)
	}
}

// TestBuildPkgTestIndexMethodCall verifies that a method call x.Foo() is
// attributed under the bare method name Foo.
func TestBuildPkgTestIndexMethodCall(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"client_test.go": `package gen
import "testing"
func TestClient(t *testing.T) {
	c := NewClient()
	c.Fetch("url")
}
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	got := lookupFuncs(idx, "Fetch")
	if len(got) != 1 || got[0] != "TestClient" {
		t.Fatalf("Fetch refs = %v, want [TestClient]", got)
	}
}

// TestBuildPkgTestIndexUntestedFunc verifies that a function no test references
// yields an empty lookup result.
func TestBuildPkgTestIndexUntestedFunc(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"only_test.go": `package gen
import "testing"
func TestOnly(t *testing.T) { Used() }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	if refs, ok := idx.Lookup("Unused"); ok || refs != nil {
		t.Fatalf("Unused should have no refs, got %v ok=%v", refs, ok)
	}
}

// TestBuildPkgTestIndexOneHopHelper verifies that a source func called only
// through a same-file non-Test helper is attributed via the 1-hop call graph.
func TestBuildPkgTestIndexOneHopHelper(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"via_helper_test.go": `package gen
import "testing"
func TestViaHelper(t *testing.T) { runScenario(t) }
func runScenario(t *testing.T) { ReadSource("http://x") }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	got := lookupFuncs(idx, "ReadSource")
	if len(got) != 1 || got[0] != "TestViaHelper" {
		t.Fatalf("ReadSource refs = %v, want [TestViaHelper] (1-hop helper)", got)
	}
	// The helper itself must also be recorded as a direct callee of the test.
	if h := lookupFuncs(idx, "runScenario"); len(h) != 1 || h[0] != "TestViaHelper" {
		t.Fatalf("runScenario refs = %v, want [TestViaHelper]", h)
	}
}

// TestBuildPkgTestIndexNoDeepIndirection verifies that indirection beyond one
// hop (test -> helper -> helper2 -> source) is NOT attributed (documented
// 1-hop boundary).
func TestBuildPkgTestIndexNoDeepIndirection(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"deep_test.go": `package gen
import "testing"
func TestDeep(t *testing.T) { hopOne(t) }
func hopOne(t *testing.T) { hopTwo(t) }
func hopTwo(t *testing.T) { DeepTarget() }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	if refs, ok := idx.Lookup("DeepTarget"); ok || refs != nil {
		t.Fatalf("DeepTarget should not be attributed past 1-hop, got %v", refs)
	}
}

// TestBuildPkgTestIndexSkipsNonTestEntries verifies that the indexer ignores
// subdirectories and non-_test.go files in the package dir, indexing only the
// _test.go files (covers the IsDir/non-test continue branch).
func TestBuildPkgTestIndexSkipsNonTestEntries(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		// A non-test source file that references Generate; it must be ignored.
		"foo.go": `package gen
func foo() { Generate("ignored", "ignored") }
`,
		// A legitimate test file; only this should be indexed.
		"generate_test.go": `package gen
import "testing"
func TestGenerate(t *testing.T) { Generate("a", "b") }
`,
	})

	// Add a subdirectory inside the package dir; it must be skipped (IsDir).
	subDir := filepath.Join(root, pkgDir, "sub")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Put a _test.go inside the subdir to prove subdirs are not descended into.
	if err := os.WriteFile(filepath.Join(subDir, "deep_test.go"), []byte(`package sub
import "testing"
func TestDeep(t *testing.T) { Buried() }
`), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	// Generate is referenced only by the indexed _test.go (TestGenerate), not
	// by foo.go (non-test source must be skipped).
	got := lookupFuncs(idx, "Generate")
	if len(got) != 1 || got[0] != "TestGenerate" {
		t.Fatalf("Generate refs = %v, want [TestGenerate]", got)
	}

	// The identifier inside the subdirectory test must not be indexed.
	if refs, ok := idx.Lookup("Buried"); ok || refs != nil {
		t.Fatalf("subdirectory _test.go must not be indexed, got %v ok=%v", refs, ok)
	}
}

// TestBuildPkgTestIndexReadDirError verifies that a non-existent package dir
// surfaces the os.ReadDir error rather than returning a usable index (covers
// the error-return branch).
func TestBuildPkgTestIndexReadDirError(t *testing.T) {
	root := t.TempDir()
	idx, err := BuildPkgTestIndex(root, filepath.Join("does", "not", "exist"))
	if err == nil {
		t.Fatal("expected error for non-existent package dir, got nil")
	}
	if idx != nil {
		t.Fatalf("expected nil index on error, got %v", idx)
	}
}

// TestMatchFuncByNameLookup verifies lookup keyed by model.Function.Name.
func TestMatchFuncByNameLookup(t *testing.T) {
	root, pkgDir := writePkg(t, map[string]string{
		"x_test.go": `package gen
import "testing"
func TestX(t *testing.T) { Generate("a", "b") }
`,
	})

	idx, err := BuildPkgTestIndex(root, pkgDir)
	if err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{Name: "Generate", QualifiedName: "internal/gen.Generate"}
	refs, ok := MatchFuncByName(idx, fn)
	if !ok || len(refs) != 1 || refs[0].TestFunc != "TestX" {
		t.Fatalf("MatchFuncByName(Generate) = %v ok=%v, want [TestX]", refs, ok)
	}

	missing := &model.Function{Name: "Nope"}
	if refs, ok := MatchFuncByName(idx, missing); ok || refs != nil {
		t.Fatalf("MatchFuncByName(Nope) = %v ok=%v, want empty", refs, ok)
	}
}
