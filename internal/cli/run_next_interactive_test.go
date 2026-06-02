package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestInteractive_unchangedPartialDoesNotTrapCursor is the ★ DoD regression
// guard: when an unchanged partial is the first TODO, calling `tsma next` must
// surface it once and then move on to the following untested TODO on the next
// call (the cursor must rotate; the partial must not re-pin itself).
func TestInteractive_unchangedPartialDoesNotTrapCursor(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	// Function 0: a partial whose test is UNCHANGED (stored mtime == current).
	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)
	partialMtime := getTestMtime(dir, filepath.Join("a", "partial_test.go"))

	// Function 1: untested (no test file) -> TODO, must become reachable.
	pb := filepath.Join(dir, "b")
	os.MkdirAll(pb, 0o755)
	os.WriteFile(filepath.Join(pb, "untested.go"),
		[]byte("package b\n\nfunc Untested() int { return 7 }\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "a.Partial", Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, CoveragePct: 50, Attempt: 1,
				TestMtime: partialMtime, TestFiles: []string{"a/partial_test.go"}},
			{QualifiedName: "b.Untested", Name: "Untested", File: "b/untested.go", StartLine: 3, EndLine: 3,
				Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 2},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Call 1: surfaces the unchanged Partial and advances the cursor.
	out1 := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("call 1 error: %v", err)
		}
	})
	if !strings.Contains(out1, "Partial") {
		t.Errorf("call 1 should surface Partial, got %q", out1)
	}
	s1 := loadSessionForTest(t, dir)
	if s1.CurrentIndex != 1 {
		t.Errorf("cursor must advance past unchanged partial to 1, got %d", s1.CurrentIndex)
	}
	if s1.Functions[0].Status != model.StatusTodo {
		t.Errorf("partial must stay TODO (no auto-DONE), got %s", s1.Functions[0].Status)
	}

	// Call 2: must now surface the following Untested TODO (not the partial).
	out2 := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("call 2 error: %v", err)
		}
	})
	if !strings.Contains(out2, "Untested") {
		t.Errorf("call 2 must reach the untested TODO, got %q", out2)
	}
}

// TestInteractive_changedPartialRemeasuredInPlace verifies the retry UX in
// interactive mode: when the partial's test is improved (changed), it is
// re-measured and the cursor stays/advances on PASS.
func TestInteractive_changedPartialRemeasuredToPass(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Already-improved test covering both branches.
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Partial(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)

	// Stored mtime is stale so detectTestChange reports changed=true.
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "partial_test.go"), future, future)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "a.Partial", Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, CoveragePct: 50, Attempt: 1,
				TestMtime: "stale", TestFiles: []string{"a/partial_test.go"}},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS after improving test, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusPass {
		t.Errorf("expected Partial -> PASS, got %s", s.Functions[0].Status)
	}
}

// TestInteractive_untestedWithMisnamedTestPrintsRename covers the rename branch
// of the untested path: when a misnamed test_<base>_test.go exists (and the
// content matcher does not attribute it), the rename instruction is surfaced and
// the cursor advances.
func TestInteractive_untestedWithMisnamedTestPrintsRename(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "thing.go"),
		[]byte("package a\n\nfunc Thing() int { return 7 }\n"), 0o644)
	// Misnamed test that does NOT reference Thing, so MatchFunc leaves tm empty
	// (falls back to the absent thing_test.go) while FindMisnamedTest detects it.
	os.WriteFile(filepath.Join(pa, "test_thing_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestUnrelated(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "a.Thing", Name: "Thing", File: "a/thing.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "test_thing_test.go") {
		t.Errorf("expected rename instruction referencing the misnamed file, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("untested function must stay TODO, got %s", s.Functions[0].Status)
	}
}

// TestInteractive_testFailKeepsTodo covers the outcomeTestFail case in
// interactive mode: a changed test that fails records the failure (status stays
// TODO) and the cursor does not move off the function.
func TestInteractive_testFailKeepsTodo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "buggy.go"),
		[]byte("package a\n\nfunc Buggy() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "buggy_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestBuggy(t *testing.T) {\n\tif Buggy() != 2 {\n\t\tt.Fatal(\"boom\")\n\t}\n}\n"), 0o644)
	// Stale stored mtime -> changed=true -> the failing test is actually run.
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "buggy_test.go"), future, future)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "a.Buggy", Name: "Buggy", File: "a/buggy.go", StartLine: 3, EndLine: 3,
				Status: model.StatusTodo, TestMtime: "stale",
				TestFiles: []string{"a/buggy_test.go"}},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "FAIL") || !strings.Contains(out, "Buggy") {
		t.Errorf("expected FAIL Buggy, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("failing function must stay TODO, got %s", s.Functions[0].Status)
	}
	if s.Functions[0].FailOutput == "" {
		t.Error("expected fail output recorded")
	}
	if s.CurrentIndex != 0 {
		t.Errorf("cursor must stay on the failing function, got %d", s.CurrentIndex)
	}
}

// TestInteractive_changedPartialSecondAttemptBecomesDone covers the outcomeDone
// case: an already-attempted partial whose changed test still leaves it <100% is
// accepted as DONE and the cursor advances.
func TestInteractive_changedPartialSecondAttemptBecomesDone(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Still only one branch covered -> stays partial.
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "partial_test.go"), future, future)

	// A trailing PASS so we can observe the cursor having advanced off Partial.
	pb := filepath.Join(dir, "b")
	os.MkdirAll(pb, 0o755)
	os.WriteFile(filepath.Join(pb, "ok.go"),
		[]byte("package b\n\nfunc Ok() int { return 1 }\n"), 0o644)

	// MaxAttempts 2 so a changed partial already at Attempt 1 reaches the
	// threshold on this (second) measured presentation -> auto-DONE.
	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		MaxAttempts:   2,
		Functions: []model.Function{
			{QualifiedName: "a.Partial", Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, CoveragePct: 50, Attempt: 1,
				TestMtime: "stale", TestFiles: []string{"a/partial_test.go"}},
			{QualifiedName: "b.Ok", Name: "Ok", File: "b/ok.go", StartLine: 3, EndLine: 3,
				Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 1, Pass: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "DONE") || !strings.Contains(out, "Partial") {
		t.Errorf("expected DONE Partial, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusDone {
		t.Errorf("expected DONE on second attempt, got %s", s.Functions[0].Status)
	}
	if s.CurrentIndex == 0 {
		t.Errorf("cursor must advance off the now-DONE function, got %d", s.CurrentIndex)
	}
}

// TestInteractive_changedPartialStaysPartial covers the outcomeRetry case in
// interactive mode: a first-attempt partial whose changed test still leaves it
// <100% keeps TODO and keeps the cursor pinned on it (so the user keeps seeing
// the function they just edited).
func TestInteractive_changedPartialStaysPartial(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Still partial (one branch).
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "partial_test.go"), future, future)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			// Attempt 0 -> attempt+1==1 (<2) -> outcomeRetry.
			{QualifiedName: "a.Partial", Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, CoveragePct: 50, Attempt: 0,
				TestMtime: "stale", TestFiles: []string{"a/partial_test.go"}},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "Partial") {
		t.Errorf("expected the just-edited Partial to be surfaced, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("partial after one edit must stay TODO, got %s", s.Functions[0].Status)
	}
	if s.Functions[0].Attempt != 1 {
		t.Errorf("attempt must increment to 1, got %d", s.Functions[0].Attempt)
	}
	if s.CurrentIndex != 0 {
		t.Errorf("cursor must stay pinned on the just-edited partial, got %d", s.CurrentIndex)
	}
}

// TestInteractive_noProgressPrintsSummary verifies that when no TODO is
// measurable, the remaining-TODO summary is shown (normal termination, not an
// error).
func TestInteractive_noProgressPrintsSummary(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg\nfunc Foo() {}"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", StartLine: 2, EndLine: 2, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "TODO function(s) remaining") {
		t.Errorf("expected remaining-TODO summary, got %q", out)
	}
}

// TestInteractive_allComplete covers the no-TODO completion path in interactive
// mode.
func TestInteractive_allComplete(t *testing.T) {
	dir := t.TempDir()
	writeGoFunc(t, dir, "A")
	sess := &model.Session{
		Project: dir, Lang: "go",
		FirstPassDone: true,
		Functions: []model.Function{
			{QualifiedName: "pkg.A", Name: "A", File: "pkg/a.go", Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Pass: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdout(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "All functions complete!") {
		t.Errorf("expected completion banner, got %q", out)
	}
}
