//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what uncoveredLocations: 미커버 브랜치의 "file:line" 위치들을 cap(10)까지 합쳐 FAIL Fact가 모델을 누락 브랜치로 곧장 가리키게 한다. cap 초과 시 "…"를 덧붙이고 멈춘다. measurement이 아니거나 Report가 nil이면 빈 문자열을 돌려준다.

package tsmagate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/reins/pkg/gate"
)

// uncoveredLocations joins the uncovered branch "file:line" locations (capped) so
// the FAIL Fact points a model straight at the missing branches.
func uncoveredLocations(ctx gate.Context) string {
	m, ok := asMeasurement(ctx)
	if !ok || m.Report == nil {
		return ""
	}
	const cap = 10
	locs := make([]string, 0, cap)
	for _, ub := range m.Report.Uncovered {
		if len(locs) == cap {
			locs = append(locs, "…")
			break
		}
		locs = append(locs, fmt.Sprintf("%s:%d", ub.File, ub.Line))
	}
	return strings.Join(locs, ", ")
}
