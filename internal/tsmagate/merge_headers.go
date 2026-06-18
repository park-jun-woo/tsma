//ff:func feature=gate type=helper control=sequence
//ff:what mergeHeaders: 누적 엔진의 헤더 합집합 어댑터(D2 1차). 기존 정명 파일의 헤더 라인들과 새 생성물의 헤더 라인들을 라인 단위 set-합집합으로 합쳐(동일 import/package 라인 dedup) 하나의 헤더 블록 라인을 돌려준다. 순서는 "기존 먼저, 새 라인 중 미출현분 뒤"로 안정적. 빈 줄·주석은 의미 중복이 없으므로 dedup 키에서 trim 후 빈 문자열은 dedup하지 않고 한 줄만 유지(가독 분리용). 복잡한 alias 충돌은 1차 비대상 — 흔한 단순 import 중복만 커버.

package tsmagate

import "strings"

// mergeHeaders unions the existing canonical header lines with the new
// generation's header lines, deduplicating identical (trimmed) header lines so a
// shared import appears once. Order is stable: existing lines first, then any new
// line not already present. Blank lines are collapsed to a single separator so
// the merged header stays readable without piling up empties.
func mergeHeaders(existing, incoming []string) []string {
	seen := map[string]bool{}
	var out []string
	add := func(lines []string) {
		for _, ln := range lines {
			key := strings.TrimSpace(ln)
			if key == "" {
				// Collapse runs of blank lines into a single separator.
				if len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
					continue
				}
				out = append(out, "")
				continue
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, ln)
		}
	}
	add(existing)
	add(incoming)
	return out
}
