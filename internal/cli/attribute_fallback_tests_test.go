package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// pyFixture writes a Python source file and its conventional test_ file, and
// returns the project root.
func pyFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "calc.py"),
		[]byte("def add(a, b):\n    return a + b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "test_calc.py"),
		[]byte("from calc import add\n\ndef test_add():\n    assert add(1, 2) == 3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestAttributeFallbackTests_pythonFileName verifies attributeFallbackTests uses
// the per-function file-name matcher for non-Go languages, attributing the
// conventional test file to a function in the matched source file.
func TestAttributeFallbackTests_pythonFileName(t *testing.T) {
	root := pyFixture(t)
	fns := []model.Function{
		{Name: "add", File: "calc.py"},
	}

	attributeFallbackTests(root, "python", fns)

	if len(fns[0].TestFiles) != 1 || filepath.Base(fns[0].TestFiles[0]) != "test_calc.py" {
		t.Errorf("add TestFiles = %v, want [test_calc.py]", fns[0].TestFiles)
	}
	if fns[0].TestFile != fns[0].TestFiles[0] {
		t.Errorf("add TestFile = %q, want %q", fns[0].TestFile, fns[0].TestFiles[0])
	}
}

// TestAttributeFallbackTests_noMatch verifies a function whose source has no
// conventional test file is left unattributed.
func TestAttributeFallbackTests_noMatch(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "lonely.py"),
		[]byte("def solo():\n    return 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fns := []model.Function{
		{Name: "solo", File: "lonely.py"},
	}

	attributeFallbackTests(root, "python", fns)

	if len(fns[0].TestFiles) != 0 || fns[0].TestFile != "" {
		t.Errorf("solo should not be attributed: TestFiles=%v TestFile=%q", fns[0].TestFiles, fns[0].TestFile)
	}
}
