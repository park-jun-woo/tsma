package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// unsafeInTest fires (LevelReview) when a matched test reaches into internals via
// the unsafe package (rulebook TS-REFL-001). It only surfaces when no Fail rule
// fired (tests pass at 100% branch coverage), turning a "covered but cheesed"
// function into a REVIEW for a human instead of a silent PASS.
var unsafeInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-001",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach into internals via unsafe",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-001"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no unsafe in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
