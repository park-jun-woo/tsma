package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// tsOwnPropertyInTest fires (LevelReview) when a matched TS test reaches private
// or non-enumerable members via Object.getOwnPropertyNames/getOwnPropertyDescriptor
// (TS-REFL-TS-003). Public iteration (Object.keys/values/entries) never fires.
var tsOwnPropertyInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-TS-003",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach internals via Object.getOwnProperty*",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-TS-003"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no Object.getOwnProperty* private reach in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
