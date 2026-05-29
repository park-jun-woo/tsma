package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
)

// captureStdoutAndStderr captures both streams produced by fn.
func captureStdoutAndStderr(fn func()) string {
	oldOut, oldErr := os.Stdout, os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	fn()

	w.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func loadSessionForTest(t *testing.T, dir string) *model.Session {
	t.Helper()
	s, err := session.Load(dir)
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	return s
}

func TestRunNext_allComplete(t *testing.T) {
	dir := t.TempDir()

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Pass: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runNext(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}

// TestRunNext_getProjectRootError covers the early getProjectRoot error branch
// (line 24) by removing the cwd so os.Getwd() fails.
func TestRunNext_getProjectRootError(t *testing.T) {
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

	if err := runNext(nil, nil); err == nil {
		t.Fatal("expected error when getProjectRoot fails")
	}
}

func TestRunNext_loadError(t *testing.T) {
	dir := t.TempDir()

	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	// Corrupt session: Exists true, Load fails -> loadOrAnalyze returns error.
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{broken"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	err := runNext(nil, nil)
	if err == nil {
		t.Fatal("expected error from corrupt session load")
	}
}

// writeGoModule writes a minimal compilable Go module with a fully covered
// function + test, returning the relative test file path.
func writeGoModule(t *testing.T, dir string) string {
	t.Helper()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Foo(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)
	return filepath.Join("pkg", "foo_test.go")
}

func writeSession(t *testing.T, dir string, sess *model.Session) {
	t.Helper()
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)
}

func TestRunNext_unchangedTestFile(t *testing.T) {
	dir := t.TempDir()
	testRel := writeGoModule(t, dir)

	// Set TestMtime to the current mtime so detectTestChange reports !changed.
	mtime := getTestMtime(dir, testRel)
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, TestMtime: mtime},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if output == "" {
		t.Error("expected non-empty output for unchanged test file")
	}
}

func TestRunNext_passOutcome(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir)

	// TestMtime 0 differs from real mtime -> changed=true -> runAndMeasure.
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(output, "PASS") {
		t.Errorf("expected PASS outcome, got: %q", output)
	}

	// Verify the session was updated and persisted to PASS.
	loaded := loadSessionForTest(t, dir)
	if loaded.Functions[0].Status != model.StatusPass {
		t.Errorf("expected Foo status PASS, got %s", loaded.Functions[0].Status)
	}
}

func TestRunNext_testFailOutcome(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	// A failing test triggers the outcomeTestFail branch.
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tt.Fatal(\"boom\")\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(output, "FAIL") {
		t.Errorf("expected FAIL outcome, got: %q", output)
	}
}

func TestRunNext_retryThenDone(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	// Function with an uncovered branch so coverage < 100%.
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	// Attempt 0 -> outcomeRetry (PARTIAL).
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out1 := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error attempt 1: %v", err)
		}
	})
	if !strings.Contains(out1, "PARTIAL") {
		t.Errorf("expected PARTIAL on first attempt, got: %q", out1)
	}

	// Touch the test file so detectTestChange reports changed on attempt 2.
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(srcDir, "foo_test.go"), future, future)

	out2 := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error attempt 2: %v", err)
		}
	})
	if !strings.Contains(out2, "DONE") {
		t.Errorf("expected DONE on second attempt, got: %q", out2)
	}

	loaded := loadSessionForTest(t, dir)
	if loaded.Functions[0].Status != model.StatusDone {
		t.Errorf("expected Foo status DONE, got %s", loaded.Functions[0].Status)
	}
}

// makeSaveFail makes session.Save fail for a project whose session is still
// loadable: it creates the .tsma/tests path as a regular file, so Save's
// os.MkdirAll(.tsma/tests) fails with "not a directory" while Load (which only
// reads session.json) still succeeds.
func makeSaveFail(t *testing.T, dir string) {
	t.Helper()
	testsPath := filepath.Join(dir, ".tsma", "tests")
	_ = os.RemoveAll(testsPath)
	if err := os.WriteFile(testsPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("seed save-fail file: %v", err)
	}
}

// TestRunNext_passThenNextTodo covers the PASS branch where advanceToNext still
// returns a following TODO function (line 82 true-branch).
func TestRunNext_passThenNextTodo(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
			{Name: "Bar", File: "pkg/bar.go", StartLine: 1, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 2},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(output, "PASS") {
		t.Errorf("expected PASS, got: %q", output)
	}
	// The next TODO (Bar) should be advertised.
	if !strings.Contains(output, "Bar") {
		t.Errorf("expected next TODO 'Bar' to be printed, got: %q", output)
	}
}

// TestRunNext_doneThenNextTodo covers the DONE branch where advanceToNext still
// returns a following TODO function (line 99 true-branch).
func TestRunNext_doneThenNextTodo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Test only covers one branch -> coverage < 100%.
	os.WriteFile(filepath.Join(srcDir, "foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	// Foo already on attempt 1 so this measurement triggers DONE, and Bar
	// remains TODO so advanceToNext returns non-nil.
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo, Attempt: 1},
			{Name: "Bar", File: "pkg/bar.go", StartLine: 1, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 2},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(output, "DONE") {
		t.Errorf("expected DONE, got: %q", output)
	}
	if !strings.Contains(output, "Bar") {
		t.Errorf("expected next TODO 'Bar' to be printed, got: %q", output)
	}
}

// TestRunNext_allCompleteSaveError covers the Save-error branch (line 36) when
// all functions are already complete.
func TestRunNext_allCompleteSaveError(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Pass: 1},
	}
	writeSession(t, dir, sess)
	makeSaveFail(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	captureStdout(func() {
		if err := runNext(nil, nil); err == nil {
			t.Fatal("expected save error when all complete")
		}
	})
}

// TestRunNext_noTestFileSaveError covers the Save-error branch (line 46) on the
// no-test-file path.
func TestRunNext_noTestFileSaveError(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg\nfunc Foo() {}"), 0o644)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 2, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)
	makeSaveFail(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	captureStdout(func() {
		if err := runNext(nil, nil); err == nil {
			t.Fatal("expected save error on no-test-file path")
		}
	})
}

// TestRunNext_unchangedSaveError covers the Save-error branch (line 55) on the
// unchanged-test-file path.
func TestRunNext_unchangedSaveError(t *testing.T) {
	dir := t.TempDir()
	testRel := writeGoModule(t, dir)
	mtime := getTestMtime(dir, testRel)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, TestMtime: mtime},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)
	makeSaveFail(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	captureStdout(func() {
		if err := runNext(nil, nil); err == nil {
			t.Fatal("expected save error on unchanged-test-file path")
		}
	})
}

// TestRunNext_finalSaveError covers the final Save-error branch (line 121)
// reached after a measurement outcome (here a PASS).
func TestRunNext_finalSaveError(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)
	makeSaveFail(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err == nil {
			t.Fatal("expected save error after measurement outcome")
		}
	})
}

// TestRunNext_misnamedTestFile covers the no-test-file path where a misnamed
// test_<base>_test.go variant exists, so the rename instruction is printed.
func TestRunNext_misnamedTestFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	// Misnamed test file (Python-convention prefix mixed in), no canonical foo_test.go.
	os.WriteFile(filepath.Join(srcDir, "test_foo_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {\n\tif Foo() != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "rename") {
		t.Errorf("expected rename instruction, got: %q", output)
	}
	if !strings.Contains(output, "test_foo_test.go") {
		t.Errorf("expected misnamed path in output, got: %q", output)
	}
	if !strings.Contains(output, filepath.Join("pkg", "foo_test.go")) {
		t.Errorf("expected canonical path in output, got: %q", output)
	}
}

func TestRunNext_todoNoTestFile(t *testing.T) {
	dir := t.TempDir()

	// Create source file but no test file
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg\nfunc Foo() {}"), 0o644)

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 2, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.MkdirAll(filepath.Join(sessDir, "tests"), 0o755)
	data, _ := json.MarshalIndent(sess, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "session.json"), data, 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	output := captureStdout(func() {
		err := runNext(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected non-empty output")
	}
}
