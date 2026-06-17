package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// rsUnsafeInTest fires (LevelReview) when a matched Rust test reaches into memory
// via an unsafe block / unsafe fn (TS-REFL-RS-001). It surfaces only when no Fail
// rule fired (tests pass at 100% branch coverage), turning a "covered but
// cheesed" function into a REVIEW instead of a silent PASS. Sibling of
// unsafeInTest (Go) / javaReflectInTest.
var rsUnsafeInTest = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "TS-REFL-RS-001",
		Level: gate.LevelReview,
		Desc:  "the matched tests do not reach into memory via unsafe",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return false, quest.Fact{}
		}
		if f := firstFinding(m, "TS-REFL-RS-001"); f != nil {
			return true, quest.Fact{Where: loc(f), Expected: "no unsafe in tests", Actual: f.Note}
		}
		return false, quest.Fact{}
	},
}
