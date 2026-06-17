package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// csReflectInfoInTest fires (LevelReview) when a matched C# test declares a
// MethodInfo/PropertyInfo/FieldInfo handle to dynamically invoke or read a
// private member (TS-REFL-CS-002). It surfaces only when no Fail rule fired
// (tests pass at 100% branch coverage), turning a "covered but cheesed" function
// into a REVIEW instead of a silent PASS. AF015 original .NET pattern.
var csReflectInfoInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-CS-002",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not declare MethodInfo/PropertyInfo/FieldInfo to reach internals",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-CS-002"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no MethodInfo/PropertyInfo/FieldInfo in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
