//ff:func feature=gate type=helper control=sequence
//ff:what asMeasurement: gate.Context.Submission에서 Prepare 결과(*measurement)를 복구한다. measurement이 아니면(방어적) ok=false를 돌려 룰이 panic 대신 "no measurement"로 크게 실패하게 한다.

package tsmagate

import "github.com/park-jun-woo/reins/pkg/gate"

// asMeasurement recovers the Prepare result from a gate.Context. A non-measurement
// submission (defensive) reports "no measurement" so a rule fails loud rather than
// panicking.
func asMeasurement(ctx gate.Context) (*measurement, bool) {
	m, ok := ctx.Submission.(*measurement)
	return m, ok
}
