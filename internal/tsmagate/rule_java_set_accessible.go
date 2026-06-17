package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// javaSetAccessibleInTest fires (LevelReview) when a matched Java test forces
// access to a private member via setAccessible(true) (TS-REFL-JV-002). It
// surfaces only when no Fail rule fired (tests pass at 100% branch coverage),
// turning a "covered but cheesed" function into a REVIEW instead of a silent
// PASS.
var javaSetAccessibleInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-JV-002",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not force private access via setAccessible(true)",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-JV-002"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no setAccessible(true) in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
