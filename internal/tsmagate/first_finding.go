//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what firstFinding: 주어진 rule ID의 첫 smell Finding을 돌려준다(없으면 nil). 한 measurement에서 룰은 smell이 반복돼도 1회만 발화한다(첫 출현이 위치를 잡고, 전수 열거는 후속 정련 — rulebook §9).

package tsmagate

import "github.com/park-jun-woo/tsma/internal/smell"

// firstFinding returns the first smell Finding with the given rule ID, or nil.
// A rule fires once per measurement even if a smell recurs (the first occurrence
// locates it; full enumeration is a later refinement, rulebook §9).
func firstFinding(m *measurement, rule string) *smell.Finding {
	for i := range m.Smells {
		if m.Smells[i].Rule == rule {
			return &m.Smells[i]
		}
	}
	return nil
}
