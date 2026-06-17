//ff:func feature=gate type=helper control=sequence
//ff:what Rules: 게이트 위반-규칙 카탈로그(rulebook §2 매핑). tests-must-pass(LevelFail, G-001)는 테스트 실행/측정 실패 시 발화하고, branch-coverage-below-100(LevelFail, G-002/G-004 — 중심 게이트)은 테스트 통과 후 AllCovered==false 시 미커버 라인을 Fact로 발화한다. 둘 다 미발화 = PASS(G-002). MaxTries→DONE(G-003)·PASS 잠금은 reins 래칫이 자동 처리하므로 룰로 만들지 않는다. reflect smell(TS-REFL-*)은 LevelReview 룰로 나중에 추가하기 쉽게 구조만 비워 둔다.

package tsmagate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Rules is the gate's violation-rule catalog. Order matters: tests-must-pass is
// listed first so that when a build is broken, gate.Evaluate selects it as
// Verdict.RootCause (the first fired Fail rule) — coverage is meaningless on a
// broken build. The coverage rule guards on !TestFailed so the two never
// double-fire on the same submission.
//
// Rulebook mapping:
//   - tests-must-pass         → G-001 (LevelFail)
//   - branch-coverage-below-100 → G-002/G-004 (LevelFail, the central gate)
//
// Not rules (reins ratchet handles them): PASS lock (G-002 success) and
// MaxTries→DONE auto-accept (G-003). The reflect-smell rules (TS-REFL-*,
// rulebook §6) are deferred; they slot in here later as LevelReview rules.
func (d *Definition) Rules() []gate.Rule {
	return []gate.Rule{testsMustPass, branchCoverageBelow100}
}

// asMeasurement recovers the Prepare result from a gate.Context. A non-measurement
// submission (defensive) reports "no measurement" so a rule fails loud rather than
// panicking.
func asMeasurement(ctx gate.Context) (*measurement, bool) {
	m, ok := ctx.Submission.(*measurement)
	return m, ok
}

// testsMustPass fires when the matched tests could not be run, did not pass, or
// the coverage tool errored (rulebook G-001). It is the upstream gate: a broken
// build short-circuits the coverage judgment.
var testsMustPass = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "tests-must-pass",
		Level: gate.LevelFail,
		Desc:  "the matched tests compile, run, and pass (and coverage measurement succeeds)",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return true, quest.Fact{Where: "submission", Expected: "a measurement", Actual: "none"}
		}
		if !m.TestFailed {
			return false, quest.Fact{}
		}
		actual := strings.TrimSpace(m.FailOutput)
		if actual == "" {
			actual = "tests did not pass"
		}
		return true, quest.Fact{
			Where:    m.FuncName,
			Expected: "tests pass",
			Actual:   firstLines(actual, 6),
		}
	},
}

// branchCoverageBelow100 is the central tsma gate (rulebook G-002/G-004). It
// fires when the tests passed but branch coverage is below 100%, carrying the
// uncovered branch locations as the Fact so a model converges on them. It is
// silent on a test failure (testsMustPass owns that) and silent at 100% (= PASS).
var branchCoverageBelow100 = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "branch-coverage-below-100",
		Level: gate.LevelFail,
		Desc:  "every branch of the function is covered (100% branch coverage)",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok || m.TestFailed || m.Report == nil {
			return false, quest.Fact{}
		}
		if m.Report.AllCovered {
			return false, quest.Fact{}
		}
		// Build a located, quantified Fact from the uncovered branches.
		where := m.FuncName
		if locs := uncoveredLocations(ctx); locs != "" {
			where = locs
		}
		return true, quest.Fact{
			Where:    where,
			Expected: "100% branch coverage",
			Actual:   fmt.Sprintf("%.1f%% (%d uncovered branch(es))", m.Report.TotalPct, len(m.Report.Uncovered)),
		}
	},
}

// uncoveredLocations joins the uncovered branch "file:line" locations (capped) so
// the FAIL Fact points a model straight at the missing branches.
func uncoveredLocations(ctx gate.Context) string {
	m, ok := asMeasurement(ctx)
	if !ok || m.Report == nil {
		return ""
	}
	const cap = 10
	locs := make([]string, 0, cap)
	for _, ub := range m.Report.Uncovered {
		if len(locs) == cap {
			locs = append(locs, "…")
			break
		}
		locs = append(locs, fmt.Sprintf("%s:%d", ub.File, ub.Line))
	}
	return strings.Join(locs, ", ")
}

// firstLines returns at most n lines of s, trimmed, so a noisy compiler dump does
// not flood the Fact.
func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append(lines[:n], "…")
	}
	return strings.Join(lines, "\n")
}
