package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestGoCheckerCheck_InvalidProject(t *testing.T) {
	checker := &GoChecker{}
	_, err := checker.Check("/nonexistent", mkMatch("fake_test.go"), nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}

// TestGoCheckerCheck_relError covers the filepath.Rel error branch (line 25):
// the test file resolves to an absolute path while projectRoot is relative, so
// the two cannot be made relative to each other.
func TestGoCheckerCheck_relError(t *testing.T) {
	checker := &GoChecker{}
	_, err := checker.Check("relative-root-not-abs", mkMatch("fake_test.go"), nil)
	if err == nil {
		t.Fatal("expected error when projectRoot is relative")
	}
	if !strings.Contains(err.Error(), "relative package") {
		t.Errorf("expected relative-package error, got: %v", err)
	}
}

// TestGoCheckerCheck_absError covers the filepath.Abs error branch (line 19):
// filepath.Abs fails only when os.Getwd() fails, which we force by removing the
// current working directory out from under the process.
func TestGoCheckerCheck_absError(t *testing.T) {
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

	checker := &GoChecker{}
	// Relative testFile so filepath.Abs must consult the (now-missing) cwd.
	if _, err := checker.Check("/proj", mkMatch("fake_test.go"), nil); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}

// writeGoCovModule builds a minimal compilable Go module with a function and a
// passing test in the pkg/ subdirectory. Returns the relative test file path.
func writeGoCovModule(t *testing.T, dir string) string {
	t.Helper()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/cov\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Foo(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)
	return filepath.Join("pkg", "foo_test.go")
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

func TestGoCheckerCheck_success(t *testing.T) {
	dir := t.TempDir()
	testRel := writeGoCovModule(t, dir)
	chdir(t, dir)

	checker := &GoChecker{}
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8}
	report, err := checker.Check(dir, mkMatch(testRel), fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if !report.AllCovered {
		t.Errorf("expected full coverage, got %+v", report)
	}
}

func TestGoCheckerCheck_testFails(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/cov\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	// Failing test makes `go test -coverprofile` exit non-zero -> error branch.
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tt.Fatal(\"boom\")\n}\n"), 0o644)
	chdir(t, dir)

	checker := &GoChecker{}
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3}
	_, err := checker.Check(dir, mkMatch(filepath.Join("pkg", "foo_test.go")), fn)
	if err == nil {
		t.Fatal("expected error when go test fails")
	}
}

// TestGoCheckerCheck_parseProfileError covers line 54: `go test` exits 0 but
// the coverprofile it produced is unparseable. A fake `go` on PATH writes a
// single line longer than bufio.Scanner's max token to the coverprofile path,
// so parseCoverProfile's scanner returns "token too long".
func TestGoCheckerCheck_parseProfileError(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/cov\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n"), 0o644)

	// Fake `go`: succeed (exit 0) but emit a >64KB single line, with no newline,
	// to whatever path is given via -coverprofile=... .
	binDir := t.TempDir()
	script := `#!/bin/sh
for a in "$@"; do
  case "$a" in
    -coverprofile=*) out="${a#-coverprofile=}" ;;
  esac
done
if [ -n "$out" ]; then
  awk 'BEGIN{ for(i=0;i<100000;i++) printf "a" }' > "$out"
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(binDir, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	checker := &GoChecker{}
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3}
	_, err := checker.Check(dir, mkMatch(filepath.Join(dir, "pkg", "foo_test.go")), fn)
	if err == nil {
		t.Fatal("expected error when coverage profile is unparseable")
	}
	if !strings.Contains(err.Error(), "parse coverage profile") {
		t.Errorf("expected 'parse coverage profile', got: %v", err)
	}
}

func TestGoCheckerCheck_partialCoverage(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/cov\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Test only exercises one branch -> report not fully covered.
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)
	chdir(t, dir)

	checker := &GoChecker{}
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8}
	report, err := checker.Check(dir, mkMatch(filepath.Join("pkg", "foo_test.go")), fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.AllCovered {
		t.Error("expected partial coverage, got fully covered")
	}
}
