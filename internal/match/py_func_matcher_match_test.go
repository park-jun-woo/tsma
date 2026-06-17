package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writePyPkg writes files (rel→content) under a fresh temp dir and returns it.
func writePyPkg(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestPythonFuncMatcher_ContentMatch proves content-aware attribution: a test
// file that calls classify() (but is NOT named test_calc.py) is still attributed
// to classify, which the filename fallback alone could never do. Skipped without
// a Python interpreter.
func TestPythonFuncMatcher_ContentMatch(t *testing.T) {
	if resolvePython() == "" {
		t.Skip("no python interpreter; content matcher unavailable")
	}
	root := writePyPkg(t, map[string]string{
		"src/calc.py": "def classify(n):\n    return 'pos' if n > 0 else 'neg'\n",
		// Deliberately NOT named test_calc.py — only a content scan can attribute it.
		"src/check_behaviour_test.py": "from calc import classify\n\ndef test_it():\n    assert classify(1) == 'pos'\n",
	})

	fn := &model.Function{Name: "classify", File: "src/calc.py"}
	tm, ok := (&PythonFuncMatcher{}).MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected a content match for classify")
	}
	if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "check_behaviour_test.py" {
		t.Fatalf("content match files = %v, want [.../check_behaviour_test.py]", tm.Files)
	}
}

// TestPythonFuncMatcher_FilenameFallback proves the last-resort fallback still
// fires when no test references the function but a conventional test_<file>.py
// exists on disk.
func TestPythonFuncMatcher_FilenameFallback(t *testing.T) {
	if resolvePython() == "" {
		t.Skip("no python interpreter; content matcher unavailable")
	}
	root := writePyPkg(t, map[string]string{
		"src/widget.py": "def render():\n    return 1\n",
		// References nothing about widget, but its name conventionally covers it.
		"src/test_widget.py": "import os\n\ndef test_env():\n    assert os.getcwd()\n",
	})

	fn := &model.Function{Name: "render", File: "src/widget.py"}
	tm, ok := (&PythonFuncMatcher{}).MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected filename fallback for render")
	}
	if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "test_widget.py" {
		t.Fatalf("fallback files = %v, want [.../test_widget.py]", tm.Files)
	}
}
