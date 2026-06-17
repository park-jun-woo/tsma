package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// csReflectInTest fires (LevelReview) when a matched C# test reaches private
// members via System.Reflection GetMethod/GetField/GetProperty
// (TS-REFL-CS-001). It surfaces only when no Fail rule fired (tests pass at 100%
// branch coverage), turning a "covered but cheesed" function into a REVIEW
// instead of a silent PASS. AF015 original .NET pattern.
var csReflectInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-CS-001",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach internals via System.Reflection GetMethod/Field/Property",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-CS-001"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no GetMethod/Field/Property in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
