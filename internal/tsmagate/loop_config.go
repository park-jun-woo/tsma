//ff:func feature=gate type=helper control=sequence
//ff:what LoopOptions: tsma의 generate→gate→retry 루프 설정값을 구성한다. DefaultModel(claude:sonnet, CLI 로그인이라 키 불필요)·System(완전한 _test.go만 출력, 펜스/산문 금지)·RuleSystem(tests-must-pass/branch-coverage-below-100 RootCause별 코칭). Escalate는 2차로 비활성(nil), EscalateOn은 능력한계 신호(커버리지 미달)만 표시해 추후 승급 대비. main이 cli.Options{Loop}에 끼운다(테스트는 LLM 필드에 stub 주입).

package tsmagate

import "github.com/park-jun-woo/reins/pkg/cli"

// LoopOptions returns the reins loop configuration for tsma: a strong default
// model (tsma targets 100% branch coverage), a system prompt that pins the output
// to a single complete _test.go file with no prose or fences, and per-RootCause
// coaching keyed to tsma's two FAIL rules. Escalate is left nil (single-backend
// first), but EscalateOn names the capability-bound RootCause (coverage still
// below 100) so wiring an Escalate backend later promotes only the residual the
// primary cannot crack. Tests inject a stub via the returned value's LLM field.
func LoopOptions() *cli.LoopOptions {
	return &cli.LoopOptions{
		// Strong default; tsma aims for 100%. claude uses CLI login (no API key).
		DefaultModel: "claude:sonnet",
		System: "You write Go tests that achieve 100% branch coverage of one function. " +
			"Output ONLY a complete, compilable _test.go file. No prose, no markdown fences.",
		RuleSystem: map[string]string{
			"tests-must-pass": "Your previous test did not compile or failed. Fix the exact error shown; " +
				"output a full corrected _test.go file.",
			"branch-coverage-below-100": "Add cases that exercise EACH uncovered file:line listed. " +
				"Cover every branch (if/else, switch, error paths, loops).",
		},
		// Escalate stays nil (single backend first); EscalateOn marks the
		// capability-bound signal so enabling Escalate later promotes only the
		// residual the primary cannot solve. A bare format slip (tests-must-pass)
		// is intentionally left off so it stays on the cheap primary.
		Escalate:   nil,
		EscalateOn: []string{"branch-coverage-below-100"},
	}
}
