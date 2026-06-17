package tsmagate

import (
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// testsMustPass fires when the matched tests could not be run, did not pass, or
// the coverage tool errored (rulebook G-001). It is the upstream gate: a broken
// build short-circuits the coverage judgment.
var testsMustPass = gate.Rule{
	Meta: gate.RuleMeta{
		ID:    "tests-must-pass",
		Level: gate.LevelFail,
		Desc:  "the matched tests compile, run, and pass (and coverage measurement succeeds)",
	},
	Check: func(ctx gate.Context) (bool, quest.Fact) {
		m, ok := asMeasurement(ctx)
		if !ok {
			return true, quest.Fact{Where: "submission", Expected: "a measurement", Actual: "none"}
		}
		if !m.TestFailed {
			return false, quest.Fact{}
		}
		if m.Truncated {
			return true, quest.Fact{
				Where:    m.FuncName,
				Expected: "complete compilable file",
				Actual:   "output appears truncated/incomplete — emit the ENTIRE file in one block",
			}
		}
		actual := strings.TrimSpace(m.FailOutput)
		if actual == "" {
			actual = "tests did not pass"
		}
		return true, quest.Fact{
			Where:    m.FuncName,
			Expected: "tests pass",
			Actual:   firstLines(actual, 6),
		}
	},
}
