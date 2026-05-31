package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writeCoveredModule writes a Go module whose test fully covers Foo.
func writeCoveredModule(t *testing.T, dir string) {
	t.Helper()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Foo(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)
}

func TestRunAndMeasure_pass(t *testing.T) {
	dir := t.TempDir()
	writeCoveredModule(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 3)

	if res.outcome != outcomePass {
		t.Fatalf("expected outcomePass, got %v (pct=%v)", res.outcome, res.coveragePct)
	}
	if res.coveragePct != 100 {
		t.Errorf("expected 100%% coverage, got %v", res.coveragePct)
	}
}

func TestRunAndMeasure_testFail(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tt.Fatal(\"boom\")\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 3)

	if res.outcome != outcomeTestFail {
		t.Fatalf("expected outcomeTestFail, got %v", res.outcome)
	}
	if res.failOutput == "" {
		t.Error("expected non-empty failOutput on test failure")
	}
}

func TestRunAndMeasure_retry(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	// Uncovered else-branch -> coverage < 100%.
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Attempt 0 -> attempt becomes 1; with maxAttempts=3 (1 < 3) -> outcomeRetry.
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Attempt: 0}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 3)

	if res.outcome != outcomeRetry {
		t.Fatalf("expected outcomeRetry, got %v", res.outcome)
	}
	if res.attempt != 1 {
		t.Errorf("expected attempt=1, got %d", res.attempt)
	}
	if len(res.uncovered) == 0 {
		t.Error("expected uncovered branches reported on retry")
	}
	if res.coveragePct >= 100 {
		t.Errorf("expected partial coverage, got %v", res.coveragePct)
	}
}

func TestRunAndMeasure_done(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Attempt 1 -> attempt becomes 2; with maxAttempts=2 (2 >= 2) -> outcomeDone.
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Attempt: 1}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 2)

	if res.outcome != outcomeDone {
		t.Fatalf("expected outcomeDone, got %v", res.outcome)
	}
	if res.attempt != 2 {
		t.Errorf("expected attempt=2, got %d", res.attempt)
	}
}

// TestRunAndMeasure_partialBelowThresholdStaysRetry verifies that a partial at
// an attempt count still below maxAttempts is reported as retry (not auto-DONE),
// exercising the configurable threshold directly.
func TestRunAndMeasure_partialBelowThresholdStaysRetry(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Attempt 1 -> attempt becomes 2; with maxAttempts=5 (2 < 5) -> outcomeRetry.
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Attempt: 1}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 5)

	if res.outcome != outcomeRetry {
		t.Fatalf("expected outcomeRetry below threshold, got %v", res.outcome)
	}
	if res.attempt != 2 {
		t.Errorf("expected attempt=2, got %d", res.attempt)
	}
}

func TestRunAndMeasure_runnerError(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Point at a test file that does not exist on disk; the Go runner fails to
	// extract test functions and returns an error -> outcomeTestFail with the
	// error string as failOutput (covers the err != nil branch).
	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "missing_test.go")), 3)

	if res.outcome != outcomeTestFail {
		t.Fatalf("expected outcomeTestFail, got %v", res.outcome)
	}
	if res.failOutput == "" {
		t.Error("expected non-empty failOutput from runner error")
	}
}

func TestRunAndMeasure_checkError(t *testing.T) {
	dir := t.TempDir()
	// Test passes (no real source), but coverage Check fails because there is
	// no valid Go package / coverable code in the resolved package path.
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	// Source file with a compile error so coverage's `go test` invocation fails
	// (the runner-level run may still classify based on its own exit, so this
	// targets the checker error path within runAndMeasure).
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\t_ = Foo()\n}\n"), 0o644)
	// Make .tsma a file so the checker cannot create the coverprofile dir,
	// forcing go test -coverprofile to fail -> checker error.
	os.WriteFile(filepath.Join(dir, ".tsma"), []byte("x"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3}
	res := runAndMeasure(dir, "go", fn, mkMatch(filepath.Join("pkg", "foo_test.go")), 3)

	if res.outcome != outcomeTestFail {
		t.Fatalf("expected outcomeTestFail from checker error, got %v", res.outcome)
	}
}
