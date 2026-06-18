//ff:func feature=gate type=test
//ff:what loop overlay 파이프라인 헬퍼 단위테스트(각 함수를 이름으로 직접 호출해 귀속): loopFailOutput(err/res/default), shouldMaterialize(실패·nil·100%·마지막시도·중간), buildLoopTestMatch(overlay 구성/쓰기실패), measureLoop(통과측정·실행실패·측정에러), promoteBacking(읽기에러·쓰기에러·성공), finalizeBacking(materialize·마지막실패폐기·보존), prepareLoopGo(잘림·빌드실패·성공 종결).

package tsmagate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

func TestLoopFailOutput_PicksSource(t *testing.T) {
	if got := loopFailOutput(errors.New("boom"), nil); got != "boom" {
		t.Errorf("err present: got %q, want boom", got)
	}
	if got := loopFailOutput(nil, &runner.Result{Output: "fail log"}); got != "fail log" {
		t.Errorf("res present: got %q, want fail log", got)
	}
	if got := loopFailOutput(nil, nil); got != "test runner returned no result" {
		t.Errorf("default: got %q", got)
	}
}

func TestShouldMaterialize_TerminalStatesOnly(t *testing.T) {
	covered := &measurement{Report: &coverage.Report{AllCovered: true}}
	partial := &measurement{Report: &coverage.Report{AllCovered: false}}
	first := &quest.Item{Tries: 0}
	last := &quest.Item{Tries: quest.MaxTries - 1}

	if shouldMaterialize(&measurement{TestFailed: true}, first) {
		t.Error("a failed measurement must never materialize")
	}
	if shouldMaterialize(&measurement{Report: nil}, first) {
		t.Error("a nil report must never materialize")
	}
	if !shouldMaterialize(covered, first) {
		t.Error("a fully-covered pass must materialize")
	}
	if !shouldMaterialize(partial, last) {
		t.Error("a pass on the final try (about to lock DONE) must materialize")
	}
	if shouldMaterialize(partial, first) {
		t.Error("a partial pass mid-loop must not materialize")
	}
}

func TestBuildLoopTestMatch_BuildsOverlay(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module bltm\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "classify.go")}}
	it := &quest.Item{Key: "pkg.Classify"}
	tm, backingRel, err := buildLoopTestMatch(p, it, classifyFullTest, []string{"TestClassify"})
	if err != nil {
		t.Fatalf("buildLoopTestMatch: %v", err)
	}
	if len(tm.Overlay) != 1 {
		t.Fatalf("expected exactly one overlay mapping, got %v", tm.Overlay)
	}
	if len(tm.Files) != 1 || !strings.HasSuffix(tm.Files[0], "zzz_tsma_gen_test.go") {
		t.Fatalf("virtual file should be the zzz placeholder in the package dir, got %v", tm.Files)
	}
	if _, statErr := os.Stat(filepath.Join(root, backingRel)); statErr != nil {
		t.Fatalf("backing file must be written under .tsma/test: %v", statErr)
	}
	// The virtual path must never reach the source package directory.
	if _, statErr := os.Stat(filepath.Join(root, tm.Files[0])); !os.IsNotExist(statErr) {
		t.Fatalf("virtual overlay file must not exist on disk, stat err = %v", statErr)
	}
}

func TestBuildLoopTestMatch_WriteError(t *testing.T) {
	root := t.TempDir()
	// Occupy .tsma with a regular file so MkdirAll(.tsma/test) fails.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "classify.go")}}
	it := &quest.Item{Key: "pkg.Classify"}
	if _, _, err := buildLoopTestMatch(p, it, classifyFullTest, []string{"TestClassify"}); err == nil {
		t.Fatal("expected a write error when .tsma cannot be created")
	}
}

func TestMeasureLoop_PassMeasures(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module mlpass\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
		"pkg/classify_test.go": classifyFullTest,
	})
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	tm := match.TestMatch{Files: []string{filepath.Join("pkg", "classify_test.go")}, TestFuncs: []string{"TestClassify"}}
	m := &measurement{}
	measureLoop(m, p, tm)
	if m.TestFailed {
		t.Fatalf("a passing full-coverage test must not fail: %s", m.FailOutput)
	}
	if m.Report == nil || !m.Report.AllCovered {
		t.Fatalf("expected a fully-covered report, got %+v", m.Report)
	}
}

func TestMeasureLoop_RunFailure(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module mlfail\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
		// A test that compiles but fails at runtime.
		"pkg/classify_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestClassify(t *testing.T) { _ = Classify(1); t.Fatal(\"boom\") }\n",
	})
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	tm := match.TestMatch{Files: []string{filepath.Join("pkg", "classify_test.go")}, TestFuncs: []string{"TestClassify"}}
	m := &measurement{}
	measureLoop(m, p, tm)
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("a failing test must set TestFailed with output, got %+v", m)
	}
}

func TestMeasureLoop_CoverageError(t *testing.T) {
	// The test runs and passes, but .tsma is occupied by a regular file so the
	// coverage step (which writes .tsma/cover.out) cannot run -> coverage error.
	root := writeGoPkg(t, map[string]string{
		"go.mod":               "module mlcov\n\ngo 1.22\n",
		"pkg/classify.go":      classifySrc,
		"pkg/classify_test.go": classifyFullTest,
	})
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	tm := match.TestMatch{Files: []string{filepath.Join("pkg", "classify_test.go")}, TestFuncs: []string{"TestClassify"}}
	m := &measurement{}
	measureLoop(m, p, tm)
	if !m.TestFailed {
		t.Fatal("a coverage-tool error after a passing run must set TestFailed")
	}
	if m.Report != nil {
		t.Fatalf("no report must be set on a coverage error, got %+v", m.Report)
	}
}

// TestMeasureLoop_VacuousZeroCoverDowngraded (C1) drives measureLoop's
// vacuous-pass guard directly: a test that compiles and exits 0 but never calls
// the target measures 0% coverage, so measureLoop downgrades it to TestFailed
// with the dedicated "exercises"/"covers 0%" Fact and leaves no Report.
func TestMeasureLoop_VacuousZeroCoverDowngraded(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module mlvac\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
		// Compiles and passes but never calls Classify -> 0% of the target.
		"pkg/classify_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestClassify(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"x\")\n\t}\n}\n",
	})
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	tm := match.TestMatch{Files: []string{filepath.Join("pkg", "classify_test.go")}, TestFuncs: []string{"TestClassify"}}
	m := &measurement{}
	measureLoop(m, p, tm)
	if !m.TestFailed {
		t.Fatal("a zero-coverage pass must be downgraded to TestFailed")
	}
	if m.Report != nil {
		t.Fatalf("a vacuous pass must carry no Report, got %+v", m.Report)
	}
	if !strings.Contains(m.FailExpected, "exercises") || !strings.Contains(m.FailOutput, "covers 0%") {
		t.Fatalf("expected the dedicated 0%% Fact, got Expected=%q Output=%q", m.FailExpected, m.FailOutput)
	}
}

// TestPrepareLoopGo_MalformedName (C1) drives prepareLoopGo's pre-measure name
// guard directly: a parseable test whose name is malformed (lowercase rune after
// "Test", which go test silently skips) is rejected before any measurement, with
// a name-specific Fact and no Report.
func TestPrepareLoopGo_MalformedName(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module plgmal\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	it := &quest.Item{Key: "pkg.Classify"}
	malformed := "package pkg\n\nimport \"testing\"\n\nfunc TestclassifyBad(t *testing.T) {\n\t_ = Classify(1)\n}\n"
	m := prepareLoopGo(it, p, []byte(malformed))
	if !m.TestFailed || m.Report != nil {
		t.Fatalf("a malformed test name must be rejected pre-measure (TestFailed, no Report), got %+v", m)
	}
	if !strings.Contains(m.FailExpected, "well-formed") || !strings.Contains(m.FailOutput, "TestclassifyBad") {
		t.Fatalf("expected a name-specific malformed Fact, got Expected=%q Output=%q", m.FailExpected, m.FailOutput)
	}
}

func TestPromoteBacking_ReadError(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module pbread\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{}
	// The backing file does not exist -> the read fails -> TestFailed.
	promoteBacking(p, m, filepath.Join(".tsma", "test", "missing.go"))
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("a missing backing must surface as TestFailed, got %+v", m)
	}
}

func TestPromoteBacking_WriteError(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module pbwrite\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	backingRel := filepath.Join(".tsma", "test", "gen.go")
	if err := writeTestFile(root, backingRel, "package pkg\n"); err != nil {
		t.Fatal(err)
	}
	// Occupy the canonical path with a directory so the write fails.
	if err := os.MkdirAll(filepath.Join(root, "pkg", "foo_test.go"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{}
	promoteBacking(p, m, backingRel)
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("a failed canonical write must surface as TestFailed, got %+v", m)
	}
}

func TestPromoteBacking_Success(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module pbok\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	backingRel := filepath.Join(".tsma", "test", "gen.go")
	body := "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { _ = Foo() }\n"
	if err := writeTestFile(root, backingRel, body); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{}
	promoteBacking(p, m, backingRel)
	if m.TestFailed {
		t.Fatalf("a clean promote must not fail: %s", m.FailOutput)
	}
	got, err := os.ReadFile(filepath.Join(root, "pkg", "foo_test.go"))
	if err != nil {
		t.Fatalf("canonical test file must exist after promote: %v", err)
	}
	// BUG-002: promote now accumulates per-function (marker-wrapped block) rather
	// than writing the raw body verbatim. The header and the test body must both
	// survive, now bounded by tsma:begin/end markers.
	gs := string(got)
	if !strings.Contains(gs, "func TestFoo(t *testing.T) { _ = Foo() }") {
		t.Fatalf("test body must survive accumulation:\n%s", gs)
	}
	if !strings.Contains(gs, "import \"testing\"") {
		t.Fatalf("header must survive accumulation:\n%s", gs)
	}
	if !strings.Contains(gs, "tsma:begin fn=") || !strings.Contains(gs, "tsma:end fn=") {
		t.Fatalf("accumulated block must carry markers:\n%s", gs)
	}
}

func TestFinalizeBacking_MaterializesOnPass(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module fbmat\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	backingRel := filepath.Join(".tsma", "test", "gen.go")
	body := "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { _ = Foo() }\n"
	if err := writeTestFile(root, backingRel, body); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{Report: &coverage.Report{AllCovered: true}}
	finalizeBacking(p, &quest.Item{Tries: 0}, m, backingRel)
	if _, err := os.Stat(filepath.Join(root, "pkg", "foo_test.go")); err != nil {
		t.Fatalf("a passing terminal result must materialize to the canonical path: %v", err)
	}
	// C2: the materialize path promotes (copies) then deletes the backing — the
	// canonical file is the real artifact, no scratch is left behind.
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Fatalf("backing must be removed after promote, stat err = %v", err)
	}
}

func TestFinalizeBacking_DiscardsOnFinalFailure(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module fbdisc\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	backingRel := filepath.Join(".tsma", "test", "gen.go")
	if err := writeTestFile(root, backingRel, "package pkg\n"); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{TestFailed: true}
	finalizeBacking(p, &quest.Item{Tries: quest.MaxTries - 1}, m, backingRel)
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Fatalf("the final-try broken backing must be discarded, stat err = %v", err)
	}
}

func TestFinalizeBacking_CleansBackingMidLoop(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module fbkeep\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	backingRel := filepath.Join(".tsma", "test", "gen.go")
	overlayRel := filepath.Join(".tsma", "test", "overlay.json")
	if err := writeTestFile(root, backingRel, "package pkg\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeTestFile(root, overlayRel, "{}"); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{File: filepath.Join("pkg", "foo.go")}}
	m := &measurement{TestFailed: true}
	// C2: a non-final failing try still sweeps the scratch — backing is
	// measurement-only and the retry Render reads only the canonical on-disk test.
	finalizeBacking(p, &quest.Item{Tries: 0}, m, backingRel)
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Fatalf("a mid-loop failing backing must be swept, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, overlayRel)); !os.IsNotExist(err) {
		t.Fatalf("the overlay JSON must be swept too, stat err = %v", err)
	}
	// And nothing was materialized to the source tree.
	if _, err := os.Stat(filepath.Join(root, "pkg", "foo_test.go")); !os.IsNotExist(err) {
		t.Fatalf("a mid-loop failure must not touch the source tree, stat err = %v", err)
	}
}

func TestPrepareLoopGo_Truncated(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module plgtr\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	it := &quest.Item{Key: "pkg.Classify"}
	truncated := "package pkg\n\nimport \"testing\"\n\nfunc TestClassify(t *testing.T) {\n\tif Classify(1) != \"pos\" {"
	m := prepareLoopGo(it, p, []byte(truncated))
	if !m.TestFailed || !m.Truncated {
		t.Fatalf("truncated source must set TestFailed+Truncated, got %+v", m)
	}
}

func TestPrepareLoopGo_BuildError(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module plgbe\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	// .tsma is a regular file -> buildLoopTestMatch's backing write fails.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	it := &quest.Item{Key: "pkg.Classify"}
	m := prepareLoopGo(it, p, []byte(classifyFullTest))
	if !m.TestFailed || m.Truncated {
		t.Fatalf("a backing write failure must be TestFailed (not truncated), got %+v", m)
	}
	if m.FailOutput == "" {
		t.Fatal("a build error must carry FailOutput")
	}
}

func TestPrepareLoopGo_SuccessMaterializes(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module plgok\n\ngo 1.22\n",
		"pkg/classify.go": classifySrc,
	})
	chdirTo(t, root)
	p := funcPayload{Lang: "go", Root: root, Fn: model.Function{
		QualifiedName: "pkg.Classify", Name: "Classify",
		File: filepath.Join("pkg", "classify.go"), StartLine: 3, EndLine: 8,
	}}
	it := &quest.Item{Key: "pkg.Classify"}
	m := prepareLoopGo(it, p, []byte(classifyFullTest))
	if m.TestFailed {
		t.Fatalf("a full covering test must pass: %s", m.FailOutput)
	}
	if m.Report == nil || !m.Report.AllCovered {
		t.Fatalf("expected a fully-covered report, got %+v", m.Report)
	}
	// A terminal pass materializes to the canonical path only — no virtual file.
	if _, err := os.Stat(filepath.Join(root, "pkg", "classify_test.go")); err != nil {
		t.Fatalf("a passing overlay measurement must materialize: %v", err)
	}
}
