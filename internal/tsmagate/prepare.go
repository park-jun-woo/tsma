//ff:func feature=gate type=helper control=sequence level=error
//ff:what Prepare: 부작용 격리 지점. payload에서 함수를 복원→(loop 모드면 raw=생성테스트를 디스크에 기록)→테스트 재매칭(match.NewFuncMatcher)→테스트 실행(runner.NewRunner.Run)→통과 시 커버리지 측정(coverage.NewChecker.Check)→결과를 measurement으로 묶어 gate.Context에 실어 반환한다. 수동 submit(MetaLoop 미설정)은 raw를 무시하고 디스크-truth만 재측정; loop(MetaLoop=true)는 정제→경로도출→쓰기 후 합류한다. 테스트 실패/측정/쓰기 에러는 measurement.TestFailed 플래그로 표현(룰이 발화하도록; short verdict 단락 대신 Rules 경유로 RootCause 보존).

package tsmagate

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// Prepare re-measures a function's tests on disk and packages the result for the
// gate. raw (the submitted bytes) is intentionally ignored: tsma's model is to
// re-measure whatever test files now exist on disk for this function, not to
// trust a submitted artifact. This is the one side-effecting step (it executes
// the test runner and coverage tool); the rules that read its result stay pure.
//
// It reconstructs the rulebook G-001/G-002/G-004 ordering from the legacy
// run_and_measure.go: a test-execution or measurement failure sets
// measurement.TestFailed (G-001) so the tests-must-pass rule fires; otherwise the
// branch-coverage Report drives the coverage gate. No short verdict is returned —
// even a broken build flows through Rules so Verdict.RootCause names the rule.
// The MaxTries→DONE auto-accept (rulebook G-003) is the reins ratchet's job, not
// a rule's, so Prepare never decides DONE.
func (d *Definition) Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	var p funcPayload
	if err := it.DecodePayload(&p); err != nil {
		return gate.Context{}, nil, fmt.Errorf("decode payload for %s: %w", it.Key, err)
	}
	fn := p.Fn

	// Native loop mode: measure the generated test non-invasively via the
	// language's own isolation — Go through `go test -overlay`, TypeScript through
	// the .tsma/test scratch with rewritten imports. The source tree is never
	// touched, so a broken generation cannot contaminate siblings and a brand-new
	// test still attributes (the disk re-match is bypassed). This returns before
	// the disk-truth path below; languages without a native path fall through.
	if isLoopMode(s) {
		if lm, ok := prepareLoopNative(it, p, raw); ok {
			return gate.Context{Item: it, Submission: lm}, nil, nil
		}
	}

	m := &measurement{FuncName: fn.QualifiedName}

	// Loop mode only (reins sets quest.MetaLoop while the `loop` command runs):
	// raw is the LLM-generated test file and must reach disk before measurement,
	// otherwise coverage never changes and every retry FAILs identically until
	// MaxTries locks DONE (improvement 0). Manual submit (MetaLoop unset) skips
	// this block, preserving the disk-truth contract — a submitted artifact is
	// never written to an arbitrary path. Order (§2-1/§2-2): derive path from a
	// pre-write match → sanitize → write; then fall through to the existing
	// re-match → run → measure so loop and submit share one measurement path. A
	// path-derivation or write failure is surfaced as TestFailed (never silent),
	// which the tests-must-pass rule turns into a FAIL the loop feeds back.
	if isLoopMode(s) {
		pre, preFound := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &fn)
		path, err := testTargetPath(p, pre, preFound)
		if err != nil {
			m.TestFailed = true
			m.FailOutput = err.Error()
			return gate.Context{Item: it, Submission: m}, nil, nil
		}
		if err := writeTestFile(p.Root, path, sanitizeGoSource(string(raw))); err != nil {
			m.TestFailed = true
			m.FailOutput = err.Error()
			return gate.Context{Item: it, Submission: m}, nil, nil
		}
	}

	// Re-match the function's tests (content-aware for Go, file-name otherwise),
	// mirroring detectTestChange so batch and single matching stay consistent. In
	// loop mode this reflects the just-written file (a new test now attributes).
	tm, found := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &fn)
	if !found || len(tm.Files) == 0 {
		// No test attributed: there is nothing to run, so this is a test failure
		// (G-001) — the function is uncovered by construction.
		m.TestFailed = true
		m.FailOutput = "no test file attributed to this function"
		return gate.Context{Item: it, Submission: m}, nil, nil
	}
	m.TestFiles = tm.Files

	// Statically scan the matched test files for escape-hatch smells (Go via
	// go/ast, TypeScript via tree-sitter — both node-based for precision). This
	// is independent of run/measure — smells are a test property, not a runtime
	// one — so it happens as soon as tm.Files is known, before the runner.
	// tm.Files are root-relative, so join with p.Root. Parse errors are ignored:
	// a broken test file is judged by tests-must-pass, not here. Findings drive
	// the LevelReview TS-REFL-* rules (surfaced only when no Fail rule fires,
	// i.e. tests pass at 100% branch coverage).
	m.Smells = scanSmells(p.Lang, p.Root, tm.Files)

	// Run the matched tests. A run error or a non-pass result is TEST_FAIL: a
	// broken build is never judged on coverage (rulebook G-001).
	res, err := runner.NewRunner(p.Lang).Run(p.Root, tm)
	if err != nil || res == nil || !res.Pass {
		m.TestFailed = true
		switch {
		case err != nil:
			m.FailOutput = err.Error()
		case res != nil:
			m.FailOutput = res.Output
		default:
			m.FailOutput = "test runner returned no result"
		}
		return gate.Context{Item: it, Submission: m}, nil, nil
	}

	// Tests pass → measure branch coverage. A measurement error is also TEST_FAIL
	// (rulebook G-001: if the coverage tool breaks, coverage is unknown).
	report, err := coverage.NewChecker(p.Lang).Check(p.Root, tm, &fn)
	if err != nil {
		m.TestFailed = true
		m.FailOutput = err.Error()
		return gate.Context{Item: it, Submission: m}, nil, nil
	}
	m.Report = report

	return gate.Context{Item: it, Submission: m}, nil, nil
}

// isLoopMode reports whether Prepare is being driven by the reins `loop` command,
// which sets quest.MetaLoop=true on the session for the loop's duration. It is
// nil-safe: a manual submit/next path (and the unit tests) may pass a nil session,
// which is treated as not-loop so the disk-truth contract holds. Only this signal
// — the run mode, not the presence of raw bytes — gates the generated-test write
// (§2-1), so submit's --in payload is never written to disk.
func isLoopMode(s *quest.Session) bool {
	if s == nil {
		return false
	}
	v, ok := s.GetMeta(quest.MetaLoop)
	return ok && v == true
}
