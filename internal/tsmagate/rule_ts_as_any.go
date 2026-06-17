package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// tsAsAnyInTest fires (LevelReview) when a matched TS test bypasses a private
// member via an `as any` cast (TS-REFL-TS-001). TS `private` is compile-time
// only, so `(x as any).priv` is the runtime escape hatch. It surfaces only when
// no Fail rule fired (tests pass at 100% branch coverage), turning a
// "covered but cheesed" function into a REVIEW instead of a silent PASS.
var tsAsAnyInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-TS-001",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not bypass private members via an `as any` cast",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-TS-001"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no `as any` private bypass in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
