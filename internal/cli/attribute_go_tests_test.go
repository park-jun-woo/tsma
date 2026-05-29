package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// goFixture writes a minimal Go module with a source file and a test file that
// references one of the functions, and returns the project root.
func goFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"),
		[]byte("module example.com/m\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "lib.go"),
		[]byte("package lib\n\nfunc Used() int { return 1 }\n\nfunc Unused() int { return 2 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "lib_test.go"),
		[]byte("package lib\n\nimport \"testing\"\n\nfunc TestUsed(t *testing.T) { _ = Used() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestAttributeGoTests_contentAware verifies attributeGoTests fills TestFiles
// for functions that a test references in its body via content-aware matching.
// Because the conventional lib_test.go exists on disk, the unreferenced Unused
// is now attributed to it via the hybrid file-name fallback (content-aware
// precision for Used, fallback safety net for Unused).
func TestAttributeGoTests_contentAware(t *testing.T) {
	root := goFixture(t)
	fns := []model.Function{
		{Name: "Used", File: "lib.go"},
		{Name: "Unused", File: "lib.go"},
	}

	attributeGoTests(root, fns)

	if len(fns[0].TestFiles) != 1 || filepath.Base(fns[0].TestFiles[0]) != "lib_test.go" {
		t.Errorf("Used TestFiles = %v, want [lib_test.go]", fns[0].TestFiles)
	}
	if fns[0].TestFile != fns[0].TestFiles[0] {
		t.Errorf("Used TestFile = %q, want %q", fns[0].TestFile, fns[0].TestFiles[0])
	}
	if len(fns[1].TestFiles) != 1 || filepath.Base(fns[1].TestFiles[0]) != "lib_test.go" {
		t.Errorf("Unused should fall back to lib_test.go: TestFiles=%v", fns[1].TestFiles)
	}
	if fns[1].TestFile != fns[1].TestFiles[0] {
		t.Errorf("Unused TestFile = %q, want %q", fns[1].TestFile, fns[1].TestFiles[0])
	}
}

// TestAttributeGoTests_emptyInput verifies an empty function slice is a no-op
// and does not panic.
func TestAttributeGoTests_emptyInput(t *testing.T) {
	attributeGoTests(t.TempDir(), nil)
}
