package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// rsTransmuteInTest fires (LevelReview) when a matched Rust test reinterprets
// bytes via std::mem::transmute (TS-REFL-RS-002). It surfaces only when no Fail
// rule fired (tests pass at 100% branch coverage).
var rsTransmuteInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-RS-002",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reinterpret memory via transmute",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-RS-002"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no transmute in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
