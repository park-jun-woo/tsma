//ff:func feature=gate type=helper control=sequence level=error
//ff:what Prepare: 부작용 격리 지점. payload에서 함수를 복원→테스트 재매칭(match.NewFuncMatcher)→테스트 실행(runner.NewRunner.Run)→통과 시 커버리지 측정(coverage.NewChecker.Check)→결과를 measurement으로 묶어 gate.Context에 실어 반환한다. 구버전 run_and_measure.go 로직을 이 시그니처로 재구성. tsma는 디스크의 테스트 파일을 재측정하는 모델이라 raw 제출 바이트는 무시한다. 테스트 실패/측정 에러는 measurement.TestFailed 플래그로 표현(룰이 발화하도록; short verdict 단락 대신 Rules 경유로 RootCause 보존).

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

	m := &measurement{FuncName: fn.QualifiedName}

	// Re-match the function's tests (content-aware for Go, file-name otherwise),
	// mirroring detectTestChange so batch and single matching stay consistent.
	tm, found := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &fn)
	if !found || len(tm.Files) == 0 {
		// No test attributed: there is nothing to run, so this is a test failure
		// (G-001) — the function is uncovered by construction.
		m.TestFailed = true
		m.FailOutput = "no test file attributed to this function"
		return gate.Context{Item: it, Submission: m}, nil, nil
	}
	m.TestFiles = tm.Files

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
