//ff:func feature=gate type=test
//ff:what Prepare 통합테스트: 디스크의 테스트를 실제 재측정한다. payload 디코드 에러·무매칭(테스트없음)·Go 100% 성공(Report 채움)·테스트 실패(TestFailed)·파싱불가 테스트 fallback(smell 스캔 continue)·비-Go(smell 스킵) 분기를 임시 fixture로 덮는다.

package tsmagate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

// chdirTo changes the working directory to dir for the duration of the test,
// restoring it on cleanup. The runner and coverage checker resolve the match's
// relative test paths against the cwd, so a Prepare run that executes tests must
// run with the project root as the cwd (as the real CLI does).
func chdirTo(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
}

func TestPrepare_DecodeError(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.TODO, Payload: json.RawMessage("{bad")}
	if _, _, err := New().Prepare(nil, it, nil, false); err == nil {
		t.Fatal("expected a decode error for invalid payload")
	}
}

func TestPrepare_NoMatch(t *testing.T) {
	// A Go function with no attributed test -> TestFailed with the
	// "no test file attributed" message (no runner is invoked).
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module preptest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, ok := asMeasurement(ctx)
	if !ok || !m.TestFailed {
		t.Fatalf("expected TestFailed measurement, got %+v (ok=%v)", m, ok)
	}
}

func TestPrepare_GoSuccessFullCoverage(t *testing.T) {
	// Fully-covered Go function: tests pass and the Report is fully covered.
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module preptest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\n" +
			"func TestFoo(t *testing.T) {\n\tif Foo() != 1 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n",
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if m.TestFailed {
		t.Fatalf("unexpected TestFailed; FailOutput=%s", m.FailOutput)
	}
	if m.Report == nil || !m.Report.AllCovered {
		t.Fatalf("expected a fully-covered report, got %+v", m.Report)
	}
}

func TestPrepare_GoTestFails(t *testing.T) {
	// A test that references Foo (so it is attributed) but fails: TestFailed.
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module preptest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\n" +
			"func TestFoo(t *testing.T) {\n\t_ = Foo()\n\tt.Fatal(\"boom\")\n}\n",
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("expected TestFailed for a failing test")
	}
}

func TestPrepare_BrokenTestFileScanSkipped(t *testing.T) {
	// The conventional foo_test.go is syntactically broken and does not reference
	// Foo, so content-aware matching misses it and the file-name fallback
	// attributes it. ScanGo then fails to parse it (the smell-scan continue
	// branch) and the runner reports a compile failure (TestFailed).
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module preptest\n\ngo 1.22\n",
		"pkg/foo.go":      "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { this is not go",
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("expected TestFailed for a broken test file")
	}
}

func TestPrepare_ManualModeIgnoresRaw(t *testing.T) {
	// Manual mode (auto=false, as reins passes for next/submit): raw is a complete
	// covering test, but the disk-truth contract means Prepare must NOT write it to
	// disk. With no test on disk the function is unattributed, so this is TestFailed
	// ("no test file attributed") — and crucially classify_test.go is never created.
	// This is the counterpart to TestPrepare_LoopModeWritesAndMeasures and pins the
	// manual/loop split that the removed quest.MetaLoop signal used to gate.
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module manualtest\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), []byte(classifyFullTest), false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "pkg", "classify_test.go")); !os.IsNotExist(statErr) {
		t.Fatalf("manual mode must not write the submitted raw to disk, stat err = %v", statErr)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("manual mode with no test on disk must be TestFailed (raw is ignored, not written)")
	}
}

func TestPrepare_LoopModeWritesAndMeasures(t *testing.T) {
	// Loop mode: the submitted raw is a complete covering test that must be
	// written to disk, then re-matched, run, and measured to full coverage.
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module looptest\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Classify", Name: "Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8}
	s := quest.New()
	ctx, _, err := New().Prepare(s, itemWithPayload(t, "go", root, fn), []byte(classifyFullTest), true)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "pkg", "classify_test.go")); statErr != nil {
		t.Fatalf("loop mode should have written the test file: %v", statErr)
	}
	m, _ := asMeasurement(ctx)
	if m.TestFailed {
		t.Fatalf("unexpected TestFailed; FailOutput=%s", m.FailOutput)
	}
	if m.Report == nil || !m.Report.AllCovered {
		t.Fatalf("expected a fully-covered report, got %+v", m.Report)
	}
}

func TestPrepare_LoopModeTargetPathError(t *testing.T) {
	// Loop mode with a non-Go function that has no test and no derivable path:
	// testTargetPath fails and Prepare surfaces it as TestFailed (never silent).
	root := t.TempDir()
	fn := model.Function{QualifiedName: "app.x", Name: "x", File: filepath.Join("app", "x.py"), StartLine: 1, EndLine: 1}
	s := quest.New()
	ctx, _, err := New().Prepare(s, itemWithPayload(t, "python", root, fn), []byte("ignored"), true)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("expected TestFailed when the test path cannot be derived")
	}
}

func TestPrepare_LoopModeMaterializeError(t *testing.T) {
	// Go loop mode where the canonical test path is occupied by a directory: the
	// overlay measurement passes at full coverage, so finalizeBacking promotes the
	// backing to the canonical path — and that write fails. promoteBacking surfaces
	// it as TestFailed (a PASS that cannot persist is fed back, never silent).
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module wfail\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	if err := os.MkdirAll(filepath.Join(root, "pkg", "foo_test.go"), 0o755); err != nil {
		t.Fatal(err)
	}
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	s := quest.New()
	full := "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { if Foo() != 1 { t.Fatal(\"x\") } }\n"
	ctx, _, err := New().Prepare(s, itemWithPayload(t, "go", root, fn), []byte(full), true)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("expected TestFailed when the canonical test file cannot be materialized")
	}
}

func TestPrepare_GoCoverageError(t *testing.T) {
	// The disk-truth path: the test runs and passes, but .tsma is occupied by a
	// regular file so the coverage step (which writes under .tsma) cannot run ->
	// a measurement error is surfaced as TestFailed (rulebook G-001).
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module covfail\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\n" +
			"func TestFoo(t *testing.T) {\n\tif Foo() != 1 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n",
	})
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed {
		t.Fatal("a coverage-tool error after a passing run must set TestFailed")
	}
	if m.Report != nil {
		t.Fatalf("no report must be set on a coverage error, got %+v", m.Report)
	}
}

func TestPrepare_RunnerErrorFromUnreadableFallback(t *testing.T) {
	// The conventional foo_test.go exists on disk as a DIRECTORY: no test
	// references Foo, so content-aware matching misses and the file-name fallback
	// attributes foo_test.go with nil TestFuncs. The runner then tries to read it
	// to extract test funcs and the read fails (it is a directory) -> runner.Run
	// returns a non-nil error, which Prepare surfaces as TestFailed (G-001).
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module runerr\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	if err := os.MkdirAll(filepath.Join(root, "pkg", "foo_test.go"), 0o755); err != nil {
		t.Fatal(err)
	}
	chdirTo(t, root)
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "go", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("a runner error must surface as TestFailed with output, got %+v", m)
	}
}

func TestPrepare_LoopModeWriteError(t *testing.T) {
	// Loop mode with a brand-new non-Go function: testTargetPath derives the
	// canonical test path, but that path is occupied by a directory so the
	// generated-test write fails -> Prepare surfaces it as TestFailed (never
	// silent). Non-Go keeps Prepare on the disk-write branch (Go uses overlay).
	root := t.TempDir()
	appDir := filepath.Join(root, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "service.py"), []byte("def do_work():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Occupy the canonical test path with a directory so the write fails.
	if err := os.MkdirAll(filepath.Join(appDir, "test_service.py"), 0o755); err != nil {
		t.Fatal(err)
	}
	fn := model.Function{QualifiedName: "app.do_work", Name: "do_work", File: filepath.Join("app", "service.py"), StartLine: 1, EndLine: 2}
	s := quest.New()
	ctx, _, err := New().Prepare(s, itemWithPayload(t, "python", root, fn), []byte("ignored"), true)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, _ := asMeasurement(ctx)
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("a generated-test write failure must surface as TestFailed, got %+v", m)
	}
}

func TestPrepare_NonGoSkipsSmellScan(t *testing.T) {
	// A non-Go (Python) function matched by file name: the smell scan (Go-only)
	// is skipped and the Python runner executes. The deliberately failing test
	// (or a missing interpreter) yields TestFailed.
	root := t.TempDir()
	appDir := filepath.Join(root, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "service.py"), []byte("def do_work():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	pyTest := "import unittest\n\nclass T(unittest.TestCase):\n    def test_x(self):\n        self.assertTrue(False)\n\nif __name__ == '__main__':\n    unittest.main()\n"
	if err := os.WriteFile(filepath.Join(appDir, "test_service.py"), []byte(pyTest), 0o644); err != nil {
		t.Fatal(err)
	}
	fn := model.Function{QualifiedName: "app.do_work", Name: "do_work", File: filepath.Join("app", "service.py"), StartLine: 1, EndLine: 2}
	ctx, _, err := New().Prepare(nil, itemWithPayload(t, "python", root, fn), nil, false)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	m, ok := asMeasurement(ctx)
	if !ok {
		t.Fatal("expected a measurement submission")
	}
	if len(m.Smells) != 0 {
		t.Errorf("expected no smells for a non-Go function, got %+v", m.Smells)
	}
	if !m.TestFailed {
		t.Fatal("expected TestFailed for the failing Python test")
	}
}
