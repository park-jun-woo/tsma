package match

import (
	"path/filepath"
	"testing"
)

// TestCanonicalPyTestPath covers the pytest same-dir naming formula and the
// non-.py rejection.
func TestCanonicalPyTestPath(t *testing.T) {
	tests := []struct {
		sourceFile string
		base       string
		want       string
	}{
		{"src/calc.py", "calc.py", filepath.Join("src", "test_calc.py")},
		{"calc.py", "calc.py", "test_calc.py"},
		{"src/notes.txt", "notes.txt", ""},
		{"src/calc.go", "calc.go", ""},
	}
	for _, tc := range tests {
		if got := canonicalPyTestPath(tc.sourceFile, tc.base); got != tc.want {
			t.Errorf("canonicalPyTestPath(%q, %q) = %q, want %q", tc.sourceFile, tc.base, got, tc.want)
		}
	}
}

// TestIsPyTestFile covers both pytest naming conventions and the rejections
// (non-.py, and a .py that is neither prefixed nor suffixed).
func TestIsPyTestFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"test_calc.py", true},
		{"calc_test.py", true},
		{"calc.py", false},
		{"test_calc.txt", false},
		{"test_data.json", false},
	}
	for _, tc := range tests {
		if got := isPyTestFile(tc.name); got != tc.want {
			t.Errorf("isPyTestFile(%q) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestPyRefsToTestMatch covers dedup (order-preserving), the empty case
// (found=false), and a normal multi-file case.
func TestPyRefsToTestMatch(t *testing.T) {
	t.Run("empty input is not found", func(t *testing.T) {
		if _, ok := pyRefsToTestMatch(nil); ok {
			t.Fatal("expected found=false for no files")
		}
	})

	t.Run("dedup preserves first-seen order", func(t *testing.T) {
		tm, ok := pyRefsToTestMatch([]string{"a_test.py", "b_test.py", "a_test.py"})
		if !ok {
			t.Fatal("expected found=true")
		}
		want := []string{"a_test.py", "b_test.py"}
		if len(tm.Files) != len(want) {
			t.Fatalf("Files = %v, want %v", tm.Files, want)
		}
		for i := range want {
			if tm.Files[i] != want[i] {
				t.Fatalf("Files = %v, want %v", tm.Files, want)
			}
		}
		if tm.TestFuncs != nil {
			t.Errorf("TestFuncs = %v, want nil (pytest runs whole files)", tm.TestFuncs)
		}
	})
}
