package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writePartialThenFullModule writes a module with a front PARTIAL function
// (uncovered branch) and a following FULLY covered function, each in its own
// package so coverage is measured independently. Returns the session functions
// in scan order: [Partial, Full].
func writePartialThenFullModule(t *testing.T, dir string) {
	t.Helper()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	// Front function: structurally <100% (one branch left uncovered).
	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	// Following function: fully covered -> should PASS even though Partial is
	// in front of it (this is the BUG-002 regression guard).
	pb := filepath.Join(dir, "b")
	os.MkdirAll(pb, 0o755)
	os.WriteFile(filepath.Join(pb, "full.go"),
		[]byte("package b\n\nfunc Full(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pb, "full_test.go"),
		[]byte("package b\n\nimport \"testing\"\n\nfunc TestFull(t *testing.T) {\n\tif Full(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Full(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)
}

// TestFirstPass_partialDoesNotBlockFollowingFull is the central DoD test: a
// front partial must not freeze CurrentIndex; the following 100% function must
// still PASS, and the partial must remain TODO (no auto-DONE).
func TestFirstPass_partialDoesNotBlockFollowingFull(t *testing.T) {
	dir := t.TempDir()
	writePartialThenFullModule(t, dir)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
			{Name: "Full", File: "b/full.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 2},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Call 1: measures Partial -> PARTIAL, watermark must advance past it.
	out1 := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("call 1 error: %v", err)
		}
	})
	if !strings.Contains(out1, "PARTIAL") || !strings.Contains(out1, "Partial") {
		t.Errorf("call 1 expected PARTIAL Partial, got %q", out1)
	}

	s1 := loadSessionForTest(t, dir)
	if s1.CurrentIndex == 0 {
		t.Errorf("watermark must advance past front partial; got CurrentIndex=%d", s1.CurrentIndex)
	}
	if s1.Functions[0].Status != model.StatusTodo {
		t.Errorf("partial must remain TODO, got %s", s1.Functions[0].Status)
	}
	if s1.Functions[0].CoveragePct == 0 {
		t.Error("partial coverage should be recorded")
	}

	// Call 2: measures Full -> PASS (the previously-blocked function).
	out2 := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("call 2 error: %v", err)
		}
	})
	if !strings.Contains(out2, "PASS") || !strings.Contains(out2, "Full") {
		t.Errorf("call 2 expected PASS Full, got %q", out2)
	}

	s2 := loadSessionForTest(t, dir)
	if s2.Functions[1].Status != model.StatusPass {
		t.Errorf("Full must PASS, got %s", s2.Functions[1].Status)
	}
	if s2.Functions[0].Status != model.StatusTodo {
		t.Errorf("partial must still be TODO after Full passes, got %s", s2.Functions[0].Status)
	}
	if !s2.FirstPassDone {
		t.Error("expected FirstPassDone after watermark reaches end")
	}
}

// TestFirstPass_allCompleteWhenWatermarkReachesEnd covers the fn==nil branch
// (advanceToNext skips every PASS and returns nil): the first pass finishes and
// the completion banner is printed. runNextFirstPass is driven directly so the
// dispatch stays in first-pass mode regardless of the FirstPassDone flag.
func TestFirstPass_allCompleteWhenWatermarkReachesEnd(t *testing.T) {
	dir := t.TempDir()
	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "A", Status: model.StatusPass, CoveragePct: 100},
			{Name: "B", Status: model.StatusDone, CoveragePct: 80},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Pass: 1, Done: 1},
	}
	writeSession(t, dir, sess)

	out := captureStdout(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "All functions complete!") {
		t.Errorf("expected completion banner, got %q", out)
	}
	if !sess.FirstPassDone {
		t.Error("expected FirstPassDone after watermark reached the end")
	}
	s := loadSessionForTest(t, dir)
	if !s.FirstPassDone {
		t.Error("expected persisted FirstPassDone")
	}
}

// TestFirstPass_untestedNoRenameAdvances covers the len(tm.Files)==0 path with
// no misnamed candidate: the function is surfaced as TODO with the generic
// next-instruction, and the watermark advances past it.
func TestFirstPass_untestedNoRenameAdvances(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "untested.go"),
		[]byte("package a\n\nfunc Untested() int { return 7 }\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Untested", File: "a/untested.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "Untested") || !strings.Contains(out, "(not found)") {
		t.Errorf("expected untested TODO with no test file, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("untested function must stay TODO, got %s", s.Functions[0].Status)
	}
	// The watermark advanced past the only function, so the first pass finished
	// (CurrentIndex is reset to 0 by finishFirstPassIfDone once it hits the end).
	if !s.FirstPassDone {
		t.Error("expected FirstPassDone after the watermark advanced past the untested function")
	}
}

// TestFirstPass_untestedWithMisnamedTestPrintsRename covers the rename branch:
// when no canonical test exists but a test_<base>_test.go is present, the rename
// instruction is surfaced instead of the generic next-instruction.
func TestFirstPass_untestedWithMisnamedTestPrintsRename(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "thing.go"),
		[]byte("package a\n\nfunc Thing() int { return 7 }\n"), 0o644)
	// Misnamed test (python-style prefix). It must NOT reference Thing, otherwise
	// the Go content-aware matcher would attribute it and tm.Files would be
	// non-empty. With no reference, MatchFunc falls back to the conventional
	// thing_test.go (absent) -> tm.Files empty, while FindMisnamedTest still
	// detects test_thing_test.go and the rename branch fires.
	os.WriteFile(filepath.Join(pa, "test_thing_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestUnrelated(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Thing", File: "a/thing.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
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
	if !s.FirstPassDone {
		t.Error("expected FirstPassDone after the watermark advanced past the function")
	}
}

// TestFirstPass_unchangedMeasuredAdvances covers the !changed branch: when the
// stored TestMtime already matches the current test mtime, the function is kept
// TODO and the watermark advances without re-measuring.
func TestFirstPass_unchangedMeasuredAdvances(t *testing.T) {
	dir := t.TempDir()
	pa := filepath.Join(dir, "a")
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	mtime := getTestMtime(dir, filepath.Join("a", "partial_test.go"))
	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, CoveragePct: 50, TestMtime: mtime,
				TestFiles: []string{"a/partial_test.go"}},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "Partial") || !strings.Contains(out, "partial_test.go") {
		t.Errorf("expected unchanged partial surfaced with its test file, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("unchanged partial must stay TODO, got %s", s.Functions[0].Status)
	}
	// Watermark advanced past the unchanged partial (the only function), so the
	// first pass finished.
	if !s.FirstPassDone {
		t.Error("expected FirstPassDone after the watermark advanced past the unchanged partial")
	}
	if s.Functions[0].CoveragePct != 50 {
		t.Errorf("coverage must be preserved (no re-measure), got %v", s.Functions[0].CoveragePct)
	}
}

// TestFirstPass_testFailDoesNotAdvance covers the outcomeTestFail case: a failing
// test pins the watermark (CurrentIndex must NOT advance) and records the failure
// output while the function stays TODO.
func TestFirstPass_testFailDoesNotAdvance(t *testing.T) {
	dir := t.TempDir()
	pa := filepath.Join(dir, "a")
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "buggy.go"),
		[]byte("package a\n\nfunc Buggy() int { return 1 }\n"), 0o644)
	// Test that fails at runtime.
	os.WriteFile(filepath.Join(pa, "buggy_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestBuggy(t *testing.T) {\n\tif Buggy() != 2 {\n\t\tt.Fatal(\"boom\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Buggy", File: "a/buggy.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "FAIL") || !strings.Contains(out, "Buggy") {
		t.Errorf("expected FAIL Buggy, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.CurrentIndex != 0 {
		t.Errorf("watermark must NOT advance past a failing test, got %d", s.CurrentIndex)
	}
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("failing function must stay TODO, got %s", s.Functions[0].Status)
	}
	if s.Functions[0].FailOutput == "" {
		t.Error("expected fail output recorded")
	}
	if s.FirstPassDone {
		t.Error("first pass must not finish while a function is failing")
	}
}

// TestFirstPass_partialOnSecondAttemptBecomesDone covers the outcomeDone case:
// a partial whose Attempt is already 1 is accepted as DONE (best-effort) on the
// next measurement, and the watermark advances.
func TestFirstPass_partialOnSecondAttemptBecomesDone(t *testing.T) {
	dir := t.TempDir()
	pa := filepath.Join(dir, "a")
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	// Test still covers only one branch -> partial.
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			// Attempt already 1 and stale mtime -> changed=true -> re-measured;
			// partial result with attempt+1>=2 yields outcomeDone.
			{Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
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
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "DONE") || !strings.Contains(out, "Partial") {
		t.Errorf("expected DONE Partial on second attempt, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusDone {
		t.Errorf("expected DONE on second attempt, got %s", s.Functions[0].Status)
	}
	// Watermark advanced past the DONE function (the only one), finishing the
	// first pass.
	if !s.FirstPassDone {
		t.Error("expected FirstPassDone after the watermark advanced past the DONE function")
	}
}

// TestFirstPass_passThenContinueSurfacesNextTodo covers the outcomePass case AND
// the post-switch continue-instruction block: a leading 100% function PASSes,
// the watermark advances, and because a following TODO remains (untested) the
// "continue" instruction plus the next TODO are surfaced. runNextFirstPass is
// driven directly so the function is attributed to it by the content-aware
// matcher.
func TestFirstPass_passThenContinueSurfacesNextTodo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	// Function 0: fully covered -> PASS.
	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "full.go"),
		[]byte("package a\n\nfunc Full(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "full_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestFull(t *testing.T) {\n\tif Full(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Full(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)

	// Function 1: untested -> stays TODO, so a TODO remains after Full passes and
	// the continue-instruction block must surface it.
	pb := filepath.Join(dir, "b")
	os.MkdirAll(pb, 0o755)
	os.WriteFile(filepath.Join(pb, "untested.go"),
		[]byte("package b\n\nfunc Untested() int { return 7 }\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Full", File: "a/full.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
			{Name: "Untested", File: "b/untested.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 2, Todo: 2},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "PASS") || !strings.Contains(out, "Full") {
		t.Errorf("expected PASS Full, got %q", out)
	}
	// continue-instruction block: the following Untested TODO is surfaced.
	if !strings.Contains(out, "Untested") {
		t.Errorf("expected the next TODO (Untested) surfaced via the continue block, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusPass {
		t.Errorf("Full must PASS, got %s", s.Functions[0].Status)
	}
	if s.CurrentIndex != 1 {
		t.Errorf("watermark must rest on the next TODO (index 1), got %d", s.CurrentIndex)
	}
	if s.FirstPassDone {
		t.Error("first pass must not be done while a TODO remains")
	}
}

// TestFirstPass_changedPartialRetries covers the outcomeRetry case: a first-pass
// function whose changed test still leaves it <100% (attempt 0 -> 1) stays TODO,
// keeps its measured coverage, and the watermark advances.
func TestFirstPass_changedPartialRetries(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)

	pa := filepath.Join(dir, "a")
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)
	// Stale stored mtime -> changed=true -> re-measured; attempt 0 -> retry.
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "partial_test.go"), future, future)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, Attempt: 0,
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
		if err := runNextFirstPass(dir, sess); err != nil {
			t.Fatalf("error: %v", err)
		}
	})
	if !strings.Contains(out, "PARTIAL") || !strings.Contains(out, "Partial") {
		t.Errorf("expected PARTIAL Partial, got %q", out)
	}
	s := loadSessionForTest(t, dir)
	if s.Functions[0].Status != model.StatusTodo {
		t.Errorf("partial must stay TODO (no auto-DONE) on first retry, got %s", s.Functions[0].Status)
	}
	if s.Functions[0].Attempt != 1 {
		t.Errorf("attempt must increment to 1, got %d", s.Functions[0].Attempt)
	}
	if s.Functions[0].CoveragePct == 0 {
		t.Error("partial coverage must be recorded")
	}
	if !s.FirstPassDone {
		t.Error("expected FirstPassDone after watermark advanced past the only function")
	}
}

// TestFirstPass_improvePartialThenPass exercises the retry UX: after the first
// pass leaves a partial TODO, improving its test (touch -> changed) re-measures
// to 100% and flips it to PASS.
func TestFirstPass_improvePartialThenPass(t *testing.T) {
	dir := t.TempDir()
	pa := filepath.Join(dir, "a")
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	os.MkdirAll(pa, 0o755)
	os.WriteFile(filepath.Join(pa, "partial.go"),
		[]byte("package a\n\nfunc Partial(x int) int {\n\tif x > 0 {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"), 0o644)
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n"), 0o644)

	sess := &model.Session{
		Project: dir, Lang: "go",
		Functions: []model.Function{
			{Name: "Partial", File: "a/partial.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Todo: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// First pass: PARTIAL, then FirstPassDone with one TODO left.
	captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("first pass error: %v", err)
		}
	})
	s1 := loadSessionForTest(t, dir)
	if !s1.FirstPassDone {
		t.Fatalf("expected FirstPassDone after single-function first pass")
	}
	if s1.Functions[0].Status != model.StatusTodo {
		t.Fatalf("partial should be TODO, got %s", s1.Functions[0].Status)
	}

	// Improve the test to cover both branches, then touch it so the change is
	// detected on the next interactive run.
	os.WriteFile(filepath.Join(pa, "partial_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestPartial(t *testing.T) {\n\tif Partial(1) != 1 {\n\t\tt.Fatal(\"a\")\n\t}\n\tif Partial(-1) != 0 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n"), 0o644)
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(pa, "partial_test.go"), future, future)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("interactive re-measure error: %v", err)
		}
	})
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS after improving the test, got %q", out)
	}
	s2 := loadSessionForTest(t, dir)
	if s2.Functions[0].Status != model.StatusPass {
		t.Errorf("expected PASS after improvement, got %s", s2.Functions[0].Status)
	}
}
