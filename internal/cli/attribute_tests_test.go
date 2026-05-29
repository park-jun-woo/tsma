package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestAttributeTests_goDispatch verifies attributeTests routes Go projects
// through the content-aware Go matcher, attributing only referenced functions.
func TestAttributeTests_goDispatch(t *testing.T) {
	root := goFixture(t)
	fns := []model.Function{
		{Name: "Used", File: "lib.go"},
		{Name: "Unused", File: "lib.go"},
	}

	attributeTests(root, "go", fns)

	if len(fns[0].TestFiles) != 1 || filepath.Base(fns[0].TestFiles[0]) != "lib_test.go" {
		t.Errorf("Used TestFiles = %v, want [lib_test.go]", fns[0].TestFiles)
	}
	if len(fns[1].TestFiles) != 0 {
		t.Errorf("Unused should not be attributed, got %v", fns[1].TestFiles)
	}
}

// TestAttributeTests_fallbackDispatch verifies attributeTests routes non-Go
// projects through the file-name fallback matcher.
func TestAttributeTests_fallbackDispatch(t *testing.T) {
	root := pyFixture(t)
	fns := []model.Function{
		{Name: "add", File: "calc.py"},
	}

	attributeTests(root, "python", fns)

	if len(fns[0].TestFiles) != 1 || filepath.Base(fns[0].TestFiles[0]) != "test_calc.py" {
		t.Errorf("add TestFiles = %v, want [test_calc.py]", fns[0].TestFiles)
	}
}
