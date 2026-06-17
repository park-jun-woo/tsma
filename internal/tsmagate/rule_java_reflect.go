package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// javaReflectInTest fires (LevelReview) when a matched Java test reaches private
// members via java.lang.reflect getDeclaredMethod/getDeclaredField
// (TS-REFL-JV-001). It surfaces only when no Fail rule fired (tests pass at 100%
// branch coverage), turning a "covered but cheesed" function into a REVIEW
// instead of a silent PASS.
var javaReflectInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-JV-001",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach internals via java.lang.reflect getDeclared*",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-JV-001"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no getDeclaredMethod/Field in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
