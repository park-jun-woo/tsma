//ff:func feature=gate type=helper control=sequence
//ff:what shouldMaterialize: overlay 측정 결과를 정명 테스트 파일로 확정(디스크 기록)할지 판단한다. 테스트 실패면 false. 통과 + 100% 커버(=PASS 확정) 또는 통과 + 이번이 마지막 시도(it.Tries==MaxTries-1, 곧 통과DONE 잠금 예정)면 true. 그 외(재시도 예정)는 false라 소스 트리에 아무것도 안 남긴다.
package tsmagate

import "github.com/park-jun-woo/reins/pkg/quest"

// shouldMaterialize reports whether the loop's overlay measurement should be
// promoted to the canonical test file on disk. A failed measurement never
// materializes. A pass at 100% coverage (the item is about to lock PASS) or a
// pass on the final try (it.Tries == MaxTries-1, about to lock DONE) does — only
// at those terminal states is the artifact committed; a mid-loop retry leaves the
// source tree untouched (overlay only).
func shouldMaterialize(m *measurement, it *quest.Item) bool {
	if m.TestFailed || m.Report == nil {
		return false
	}
	return m.Report.AllCovered || it.Tries == quest.MaxTries-1
}
