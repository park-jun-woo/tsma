//ff:type feature=gate type=model
//ff:what tsma의 reins gate.Definition 구현체와 공유 타입. Definition은 Seed/Render/Prepare/Rules 4메서드로 reins 래칫에 도메인 로직(인덱싱·매칭·테스트실행·커버리지측정·100%브랜치 게이트)을 끼운다. funcPayload는 Item.Payload에 담는 도메인 스냅샷(lang/root/함수 위치; Status/Attempt는 reins Item.State/Tries가 소유하므로 제외), measurement는 Prepare가 디스크를 측정해 만든 결과를 Rules로 운반하는 Context.Submission 타입.

package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

// Definition implements gate.Definition for tsma. It carries no mutable state;
// per-quest data rides in each Item's payload, so a single value is shared
// across all subcommands.
type Definition struct{}

// New returns a tsma gate Definition.
func New() *Definition { return &Definition{} }

// Ensure Definition satisfies the reins contract at compile time.
var _ gate.Definition = (*Definition)(nil)

// funcPayload is the domain snapshot stored in quest.Item.Payload. It carries
// everything Render/Prepare need to re-locate, re-match, and re-measure a
// function deterministically. lang and root are duplicated into every item
// because Seed has no session to stash them in (reins seeds items, not Meta),
// keeping each item self-contained. Fn deliberately omits the runtime fields
// reins now owns: Status (→ Item.State) and Attempt (→ Item.Tries) are zeroed at
// seed time so the ratchet is the single source of truth for progress.
type funcPayload struct {
	Lang string         `json:"lang"`
	Root string         `json:"root"`
	Fn   model.Function `json:"fn"`
}

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
	// Report is the branch-coverage result; nil when TestFailed.
	Report *coverage.Report
	// FuncName is the qualified name of the function under test (for Fact.Where).
	FuncName string
	// TestFiles are the matched test files (for Fact.Where when no test matched).
	TestFiles []string
}
