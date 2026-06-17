package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// rsPtrInTest fires (LevelReview) when a matched Rust test forces raw-pointer
// access via std::ptr read/write or as_ptr (TS-REFL-RS-003). It surfaces only
// when no Fail rule fired (tests pass at 100% branch coverage).
var rsPtrInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-RS-003",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not force raw-pointer access via std::ptr",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-RS-003"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no std::ptr raw access in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
