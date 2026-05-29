package coverage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestPyCheckerCheck_InvalidProject(t *testing.T) {
	checker := &PyChecker{}
	_, err := checker.Check("/nonexistent", mkMatch("fake_test.py"), nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}

// TestPyCheckerCheck_absError covers the filepath.Abs error branch (line 16):
// Abs fails only when os.Getwd() fails, forced by removing the cwd.
func TestPyCheckerCheck_absError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	checker := &PyChecker{}
	if _, err := checker.Check("/proj", mkMatch("rel_test.py"), nil); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}

// TestPyCheckerCheck_success covers the parse + report success path (lines 35,
// 36, 40) by injecting a fake `python3` that emulates pytest-cov: it parses the
// --cov-report=json:PATH argument and writes a valid coverage.json there.
func TestPyCheckerCheck_success(t *testing.T) {
	proj := t.TempDir()
	os.WriteFile(filepath.Join(proj, "mod.py"),
		[]byte("def foo(x):\n    if x:\n        return 1\n    return 0\n"), 0o644)
	os.WriteFile(filepath.Join(proj, "mod_test.py"),
		[]byte("def test_foo():\n    assert True\n"), 0o644)

	binDir := t.TempDir()
	// The fake python writes a coverage.json to the path given in the
	// --cov-report=json:<path> argument, then exits 0. The JSON references
	// mod.py with executed and missing lines so buildPyReport has data.
	script := `#!/bin/sh
for arg in "$@"; do
  case "$arg" in
    --cov-report=json:*)
      out="${arg#--cov-report=json:}"
      cat > "$out" <<'JSON'
{"files":{"mod.py":{"executed_lines":[1,2,3],"missing_lines":[4]}}}
JSON
      ;;
  esac
done
exit 0
`
	if err := os.WriteFile(filepath.Join(binDir, "python3"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	checker := &PyChecker{}
	fn := &model.Function{Name: "foo", File: "mod.py", StartLine: 1, EndLine: 4}
	report, err := checker.Check(proj, mkMatch("mod_test.py"), fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
}

// TestPyCheckerCheck_parseFails covers the parse-coverage-json error branch
// (lines 35-37): pytest-cov "succeeds" but writes no coverage.json, so the
// subsequent parse fails.
func TestPyCheckerCheck_parseFails(t *testing.T) {
	proj := t.TempDir()
	os.WriteFile(filepath.Join(proj, "mod_test.py"), []byte("def test_x():\n    assert True\n"), 0o644)

	binDir := t.TempDir()
	// Fake python that exits 0 but never writes the JSON file.
	if err := os.WriteFile(filepath.Join(binDir, "python3"), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	checker := &PyChecker{}
	fn := &model.Function{Name: "x", File: "mod_test.py", StartLine: 1, EndLine: 2}
	_, err := checker.Check(proj, mkMatch("mod_test.py"), fn)
	if err == nil {
		t.Fatal("expected parse error when coverage.json missing")
	}
}

func TestPyCheckerCheck_mkdirError(t *testing.T) {
	dir := t.TempDir()
	// Make .tsma a regular file so MkdirAll(.tsma) fails.
	if err := os.WriteFile(filepath.Join(dir, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	checker := &PyChecker{}
	fn := &model.Function{Name: "foo", File: "mod.py", StartLine: 1, EndLine: 3}
	_, err := checker.Check(dir, mkMatch("mod.py"), fn)
	if err == nil {
		t.Fatal("expected error when .tsma directory cannot be created")
	}
}

func TestPyCheckerCheck_coverageFails(t *testing.T) {
	dir := t.TempDir()
	// A test file that fails to import/run, so both pytest-cov and the
	// coverage.py fallback error out -> "python coverage failed" branch.
	os.WriteFile(filepath.Join(dir, "broken_test.py"),
		[]byte("import this_module_does_not_exist_xyz\n\ndef test_x():\n    assert False\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	checker := &PyChecker{}
	fn := &model.Function{Name: "x", File: "broken_test.py", StartLine: 1, EndLine: 4}
	_, err := checker.Check(dir, mkMatch("broken_test.py"), fn)
	if err == nil {
		t.Fatal("expected error when python coverage cannot run")
	}
}
