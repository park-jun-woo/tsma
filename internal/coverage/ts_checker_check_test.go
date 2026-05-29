package coverage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestTSCheckerCheck_InvalidProject(t *testing.T) {
	checker := &TSChecker{}
	_, err := checker.Check("/nonexistent", mkMatch("fake_test.ts"), nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}

// TestTSCheckerCheck_absError covers the filepath.Abs error branch (line 17):
// Abs fails only when os.Getwd() fails, forced by removing the cwd.
func TestTSCheckerCheck_absError(t *testing.T) {
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

	checker := &TSChecker{}
	if _, err := checker.Check("/proj", mkMatch("rel.test.ts"), nil); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}

// TestTSCheckerCheck_relFallback covers the filepath.Rel error fallback (line
// 22: relTest = absTest). When projectRoot is relative, Rel(relativeRoot,
// absTest) cannot make the two relative, so the fallback assigns absTest.
// A fake npx that exits non-zero makes the call still return an error after
// the fallback line executes.
func TestTSCheckerCheck_relFallback(t *testing.T) {
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "npx"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	// Run from a writable temp dir so the relative projectRoot's .tsma dir can
	// be created.
	work := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(work); err != nil {
		t.Fatal(err)
	}

	checker := &TSChecker{}
	fn := &model.Function{Name: "App", File: "src/app.ts", StartLine: 1, EndLine: 10}
	// Relative projectRoot triggers the Rel-error fallback at line 22.
	_, err := checker.Check("relroot", mkMatch("src/app.test.ts"), fn)
	if err == nil {
		t.Fatal("expected error (npx exits non-zero) after rel fallback")
	}
}

// fakeNpx installs a fake `npx` on PATH with the supplied shell body and
// chdirs into a fresh project dir (so filepath.Abs/Rel resolve correctly).
// Returns the project dir.
func fakeNpx(t *testing.T, body string) string {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "npx"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	proj := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(proj); err != nil {
		t.Fatal(err)
	}
	return proj
}

func TestTSCheckerCheck_success(t *testing.T) {
	// Fake npx writes a valid istanbul coverage-final.json into the coverage dir
	// (cmd.Dir is the project root, so the relative path resolves there).
	proj := fakeNpx(t, `
mkdir -p .tsma/coverage
cat > .tsma/coverage/coverage-final.json <<'JSON'
{"src/app.ts":{"path":"src/app.ts","statementMap":{},"s":{},"branchMap":{},"b":{},"fnMap":{},"f":{}}}
JSON
exit 0
`)

	checker := &TSChecker{}
	fn := &model.Function{Name: "App", File: "src/app.ts", StartLine: 1, EndLine: 10}
	report, err := checker.Check(proj, mkMatch("src/app.test.ts"), fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
}

func TestTSCheckerCheck_npxFails(t *testing.T) {
	proj := fakeNpx(t, "echo boom 1>&2\nexit 1\n")

	checker := &TSChecker{}
	fn := &model.Function{Name: "App", File: "src/app.ts", StartLine: 1, EndLine: 10}
	_, err := checker.Check(proj, mkMatch("src/app.test.ts"), fn)
	if err == nil {
		t.Fatal("expected error when npx coverage fails")
	}
}

func TestTSCheckerCheck_parseFails(t *testing.T) {
	// npx succeeds but produces no coverage-final.json -> parse error.
	proj := fakeNpx(t, "exit 0\n")

	checker := &TSChecker{}
	fn := &model.Function{Name: "App", File: "src/app.ts", StartLine: 1, EndLine: 10}
	_, err := checker.Check(proj, mkMatch("src/app.test.ts"), fn)
	if err == nil {
		t.Fatal("expected parse error when coverage report is missing")
	}
}

func TestTSCheckerCheck_mkdirFails(t *testing.T) {
	proj := fakeNpx(t, "exit 0\n")
	// Make .tsma a file so MkdirAll(.tsma/coverage) fails.
	if err := os.WriteFile(filepath.Join(proj, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	checker := &TSChecker{}
	fn := &model.Function{Name: "App", File: "src/app.ts", StartLine: 1, EndLine: 10}
	_, err := checker.Check(proj, mkMatch("src/app.test.ts"), fn)
	if err == nil {
		t.Fatal("expected error when coverage dir cannot be created")
	}
}
