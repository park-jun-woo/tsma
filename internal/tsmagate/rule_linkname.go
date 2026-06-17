package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// linknameInTest fires (LevelReview) when a matched test pulls another package's
// private symbols via a //go:linkname directive (TS-REFL-003).
var linknameInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-003",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not pull private symbols via //go:linkname",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-003"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no //go:linkname in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
