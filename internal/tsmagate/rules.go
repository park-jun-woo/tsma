//ff:func feature=gate type=helper control=sequence
//ff:what Rules: 게이트 위반-규칙 카탈로그(rulebook §2 매핑). tests-must-pass(LevelFail, G-001)는 테스트 실행/측정 실패 시 발화하고, branch-coverage-below-100(LevelFail, G-002/G-004 — 중심 게이트)은 테스트 통과 후 AllCovered==false 시 미커버 라인을 Fact로 발화한다. 그 뒤 TS-REFL-001/002/003(LevelReview, rulebook §6)은 measurement.Smells에 해당 escape-hatch Finding이 있으면 발화한다. Fail 룰이 앞이라 커버리지 FAIL이면 smell은 묻히고(정상), 100% PASS인데 smell이 있으면 REVIEW가 된다. 둘 다·셋 다 미발화 = PASS(G-002). MaxTries→DONE(G-003)·PASS 잠금은 reins 래칫이 자동 처리하므로 룰로 만들지 않는다.

package tsmagate

import "github.com/park-jun-woo/reins/pkg/gate"

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
// Rulebook §6 mapping (LevelReview, appended after the Fail rules so a coverage
// FAIL wins RootCause and a smell only surfaces as REVIEW when tests pass at
// 100%):
//   - TS-REFL-001 → unsafe-in-test
//   - TS-REFL-002 → reflect-dynamic-in-test
//   - TS-REFL-003 → linkname-in-test
//
// Not rules (reins ratchet handles them): PASS lock (G-002 success) and
// MaxTries→DONE auto-accept (G-003).
func (d *Definition) Rules() []gate.Rule {
	return []gate.Rule{
		testsMustPass,
		branchCoverageBelow100,
		unsafeInTest,
		reflectDynamicInTest,
		linknameInTest,
		tsAsAnyInTest,
		tsReflectInTest,
		tsOwnPropertyInTest,
		javaReflectInTest,
		javaSetAccessibleInTest,
		csReflectInTest,
		csReflectInfoInTest,
		rsUnsafeInTest,
		rsTransmuteInTest,
		rsPtrInTest,
	}
}
