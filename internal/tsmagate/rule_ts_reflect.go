package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// tsReflectInTest fires (LevelReview) when a matched TS test introspects via the
// Reflect API (TS-REFL-TS-002), the dynamic private-reach escape hatch.
var tsReflectInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-TS-002",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach internals via the Reflect API",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-TS-002"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no Reflect.* in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
