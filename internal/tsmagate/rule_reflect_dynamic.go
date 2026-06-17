package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// reflectDynamicInTest fires (LevelReview) when a matched test reaches private
// methods/fields dynamically via reflect MethodByName/FieldByName (TS-REFL-002).
// It never fires on reflect.DeepEqual/TypeOf/ValueOf (the detector only matches
// those two selector names).
var reflectDynamicInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-002",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not dynamically reach internals via reflect MethodByName/FieldByName",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-002"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no reflect MethodByName/FieldByName in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
