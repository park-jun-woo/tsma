package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writePyFiles writes rel→content under a fresh temp dir and returns it.
func writePyFiles(t *testing.T, files map[string]string) string {
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

// TestCollectPyCalledNames_Success runs the real ast names script and confirms
// it returns referenced names. Skipped without a Python interpreter.
func TestCollectPyCalledNames_Success(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "test_calc.py")
	src := "from calc import classify\n\ndef test_it():\n    assert classify(1)\n"
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	names := collectPyCalledNames(python, f)
	got := map[string]bool{}
	for _, n := range names {
		got[n] = true
	}
	if !got["classify"] {
		t.Errorf("names = %v, want to contain classify", names)
	}
}

// TestCollectPyCalledNames_RunError proves an unrunnable interpreter yields nil
// (one bad file never aborts the index).
func TestCollectPyCalledNames_RunError(t *testing.T) {
	if names := collectPyCalledNames("definitely-not-a-real-interpreter-xyz", "x.py"); names != nil {
		t.Fatalf("names = %v, want nil on interpreter failure", names)
	}
}

// TestCollectPyCalledNames_JSONError covers the json-unmarshal failure branch: a
// stub interpreter that exits 0 but prints non-JSON.
func TestCollectPyCalledNames_JSONError(t *testing.T) {
	dir := t.TempDir()
	stub := filepath.Join(dir, "fakepy")
	if err := os.WriteFile(stub, []byte("#!/bin/sh\necho 'not json at all'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if names := collectPyCalledNames(stub, "x.py"); names != nil {
		t.Fatalf("names = %v, want nil on non-JSON output", names)
	}
}

// TestIngestPyDir covers a missing directory (no-op), and a present directory
// where only the test files (not subdirs or non-test files) are ingested.
func TestIngestPyDir(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}

	t.Run("missing directory is a no-op", func(t *testing.T) {
		idx := &PyPkgTestIndex{refs: map[string][]string{}}
		ingestPyDir(idx, python, t.TempDir(), "does/not/exist")
		if len(idx.refs) != 0 {
			t.Fatalf("refs = %v, want empty for a missing dir", idx.refs)
		}
	})

	t.Run("only test files are ingested", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"pkg/test_calc.py": "from calc import classify\n\ndef test_x():\n    assert classify(1)\n",
			"pkg/calc.py":      "def classify(n):\n    return n\n",
			"pkg/sub/test_other.py": "from other import widget\n\ndef test_y():\n    assert widget()\n",
		})
		idx := &PyPkgTestIndex{refs: map[string][]string{}}
		ingestPyDir(idx, python, root, "pkg")
		if _, ok := idx.refs["classify"]; !ok {
			t.Errorf("expected classify ref from pkg/test_calc.py; got %v", idx.refs)
		}
		if _, ok := idx.refs["widget"]; ok {
			t.Errorf("subdir test file should not be ingested by ingestPyDir; got %v", idx.refs)
		}
	})
}

// TestIngestPyTestFile records a back-reference for every referenced name.
func TestIngestPyTestFile(t *testing.T) {
	python := resolvePython()
	if python == "" {
		t.Skip("no python interpreter on PATH")
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "test_calc.py")
	if err := os.WriteFile(f, []byte("from calc import classify\n\ndef test_it():\n    assert classify(1)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := &PyPkgTestIndex{refs: map[string][]string{}}
	ingestPyTestFile(idx, python, f, "rel/test_calc.py")
	files := idx.refs["classify"]
	if len(files) != 1 || files[0] != "rel/test_calc.py" {
		t.Fatalf("refs[classify] = %v, want [rel/test_calc.py]", files)
	}
}

// TestBuildPyPkgTestIndex covers the no-interpreter nil, the empty-index nil, and
// the populated index branches.
func TestBuildPyPkgTestIndex(t *testing.T) {
	t.Run("nil when no interpreter", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if idx := BuildPyPkgTestIndex(t.TempDir(), "pkg"); idx != nil {
			t.Fatalf("idx = %+v, want nil without an interpreter", idx)
		}
	})

	if resolvePython() == "" {
		t.Skip("no python interpreter on PATH for the populated cases")
	}

	t.Run("nil when nothing referenced", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"pkg/calc.py": "def classify(n):\n    return n\n",
		})
		if idx := BuildPyPkgTestIndex(root, "pkg"); idx != nil {
			t.Fatalf("idx = %+v, want nil when no test files exist", idx)
		}
	})

	t.Run("populated when a test references a name", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"pkg/calc.py":      "def classify(n):\n    return n\n",
			"pkg/test_calc.py": "from calc import classify\n\ndef test_it():\n    assert classify(1)\n",
		})
		idx := BuildPyPkgTestIndex(root, "pkg")
		if idx == nil {
			t.Fatal("expected a populated index")
		}
		if _, ok := idx.refs["classify"]; !ok {
			t.Errorf("index missing classify ref; got %v", idx.refs)
		}
	})
}

// TestPyFilenameFallback covers both arms: a conventional test file present
// (found) and absent (not found).
func TestPyFilenameFallback(t *testing.T) {
	t.Run("found when test_<file>.py exists", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"src/widget.py":      "def render():\n    return 1\n",
			"src/test_widget.py": "def test_env():\n    assert True\n",
		})
		fn := &model.Function{Name: "render", File: "src/widget.py"}
		tm, ok := pyFilenameFallback(root, fn)
		if !ok {
			t.Fatal("expected a filename fallback match")
		}
		if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "test_widget.py" {
			t.Fatalf("Files = %v, want [.../test_widget.py]", tm.Files)
		}
	})

	t.Run("not found when no conventional test exists", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"src/widget.py": "def render():\n    return 1\n",
		})
		fn := &model.Function{Name: "render", File: "src/widget.py"}
		if _, ok := pyFilenameFallback(root, fn); ok {
			t.Fatal("expected no fallback when no test file exists")
		}
	})
}

// TestPythonFuncMatcher_NilAndNoInterpreter covers MatchFunc's nil guard and the
// idx==nil (no interpreter) path that delegates to the filename fallback.
func TestPythonFuncMatcher_NilAndNoInterpreter(t *testing.T) {
	t.Run("nil function is not matched", func(t *testing.T) {
		if _, ok := (&PythonFuncMatcher{}).MatchFunc(t.TempDir(), nil); ok {
			t.Fatal("expected found=false for a nil function")
		}
	})

	t.Run("no interpreter falls back to filename match", func(t *testing.T) {
		root := writePyFiles(t, map[string]string{
			"src/widget.py":      "def render():\n    return 1\n",
			"src/test_widget.py": "def test_env():\n    assert True\n",
		})
		t.Setenv("PATH", t.TempDir()) // force resolvePython()=="" → idx==nil
		fn := &model.Function{Name: "render", File: "src/widget.py"}
		tm, ok := (&PythonFuncMatcher{}).MatchFunc(root, fn)
		if !ok {
			t.Fatal("expected filename fallback when no interpreter is present")
		}
		if len(tm.Files) != 1 || filepath.Base(tm.Files[0]) != "test_widget.py" {
			t.Fatalf("Files = %v, want [.../test_widget.py]", tm.Files)
		}
	})
}
