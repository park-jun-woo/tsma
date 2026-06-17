//ff:func feature=gate type=helper control=sequence level=error
//ff:what measureLoop: overlay TestMatch로 테스트를 실행(runner.Run)→통과 시 커버리지 측정(coverage.Check)해 measurement에 결과를 채운다. runner/checker는 TestMatch.Overlay가 있으면 -overlay -vet=off로 가상 테스트를 패키지에 끼워 측정한다(소스 트리 무침습). 실행 실패/측정 에러는 TestFailed로(G-001). C1 vacuous-pass 가드: 실행은 통과인데 대상 함수 커버리지가 0%면(무검증 통과) TestFailed로 강등하고 전용 Fact를 실어 materialize를 막는다. 디스크 재측정 경로(prepare.go)와 같은 형태지만 MatchFunc를 안 타고 overlay TestMatch를 그대로 받는다.
package tsmagate

import (
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// measureLoop runs the overlay TestMatch and, on pass, measures branch coverage,
// filling m. The Go runner/checker honor tm.Overlay (-overlay -vet=off) so the
// generated test is spliced into the package without touching the source tree. A
// run failure or a measurement error is TestFailed (rulebook G-001), mirroring
// the disk path's shape but taking the pre-built overlay TestMatch directly
// instead of re-matching on disk.
func measureLoop(m *measurement, p funcPayload, tm match.TestMatch) {
	res, err := runner.NewRunner(p.Lang).Run(p.Root, tm)
	if err != nil || res == nil || !res.Pass {
		m.TestFailed, m.FailOutput = true, loopFailOutput(err, res)
		return
	}
	report, err := coverage.NewChecker(p.Lang).Check(p.Root, tm, &p.Fn)
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return
	}
	// C1 vacuous-pass guard: a test that compiles and exits 0 but covers 0% of the
	// target did not exercise it (a malformed name go test ignored, a t.Skip, an
	// empty body, or a test of the wrong function). Downgrade to TestFailed with a
	// dedicated Fact so (a) tests-must-pass speaks, (b) shouldMaterialize is false,
	// and (c) finalizeBacking's final-try discard path absorbs it — no separate
	// branch. The 0 compare is via a tiny epsilon (any real call covers ≫0%).
	if report != nil && report.TotalPct < 1e-9 {
		m.TestFailed = true
		m.FailExpected = "a test that exercises " + p.Fn.QualifiedName
		m.FailOutput = "generated test runs but covers 0% of " + p.Fn.QualifiedName +
			" (does it call it? is the test name well-formed: TestXxx with an uppercase letter after 'Test'?)"
		return
	}
	m.Report = report
}
