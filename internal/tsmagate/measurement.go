//ff:type feature=gate type=model
//ff:what measurement은 Prepare가 디스크를 재측정해 만든 결과를 gate.Context.Submission으로 Rules에 운반하는 타입이다. tsma는 디스크의 테스트 파일을 재측정하므로 raw 제출 바이트는 무시한다. 두 실패 형태 중 정확히 하나가 세팅된다: TestFailed(테스트 미통과/커버리지 측정 에러 — G-001)는 FailOutput을, 아니면 Report가 브랜치커버리지 결과(G-002/G-004)를 담는다. Smells는 매칭된 테스트의 escape-hatch Finding으로 LevelReview TS-REFL-* 룰을 구동한다.

package tsmagate

import (
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/smell"
)

// measurement is the decoded submission a Prepare run hands to the gate rules
// via gate.Context.Submission. It is built from a disk re-measurement (tsma's
// model is to re-measure the test files on disk, so the raw submit bytes are
// ignored). Exactly one of the two failure shapes is set: TestFailed (tests did
// not compile/pass, or coverage measurement errored — rulebook G-001) carries
// FailOutput; otherwise Report holds the branch-coverage result (G-002/G-004).
type measurement struct {
	// TestFailed is true when the matched tests could not be run, did not pass,
	// or the coverage tool errored. It short-circuits the coverage gate so a
	// broken build is never judged on coverage.
	TestFailed bool
	// FailOutput is the test/measurement failure text surfaced as feedback.
	FailOutput string
	// Truncated is true when (loop mode, C3) the generated source did not parse
	// after sanitize+tidy — the model's output was cut off. It implies TestFailed
	// and makes tests-must-pass emit a dedicated "emit the ENTIRE file" Fact
	// instead of a raw compiler error, so a weak model can self-correct.
	Truncated bool
	// Report is the branch-coverage result; nil when TestFailed.
	Report *coverage.Report
	// FuncName is the qualified name of the function under test (for Fact.Where).
	FuncName string
	// TestFiles are the matched test files (for Fact.Where when no test matched).
	TestFiles []string
	// Smells are the test-smell findings (escape hatches) detected by scanning
	// the matched test files statically. Independent of TestFailed/Report: it is
	// populated from the test sources once they are matched, and drives the
	// LevelReview TS-REFL-* rules. Empty when no smell was found (the clean case).
	Smells []smell.Finding
}
