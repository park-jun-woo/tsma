package tsmagate

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// branchCoverageBelow100 is the central tsma gate (rulebook G-002/G-004). It
// fires when the tests passed but branch coverage is below 100%, carrying the
// uncovered branch locations as the Fact so a model converges on them. It is
// silent on a test failure (testsMustPass owns that) and silent at 100% (= PASS).
var branchCoverageBelow100 = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "branch-coverage-below-100",
		Level: gate.LevelFail,
		Desc:  "every branch of the function is covered (100% branch coverage)",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok || m.TestFailed || m.Report == nil {
			return false, quest.Fact{}
		}
		if m.Report.AllCovered {
			return false, quest.Fact{}
		}
		// Build a located, quantified Fact from the uncovered branches.
		where := m.FuncName
		if locs := uncoveredLocations(ctx); locs != "" {
			where = locs
		}
		return true, quest.Fact{
			Where:    where,
			Expected: "100% branch coverage",
			Actual:   fmt.Sprintf("%.1f%% (%d uncovered branch(es))", m.Report.TotalPct, len(m.Report.Uncovered)),
		}
	},
}
