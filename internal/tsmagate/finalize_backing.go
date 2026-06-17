//ff:func feature=gate type=helper control=sequence lang=go
//ff:what finalizeBacking: overlay 측정 종결 처리(C2). 확정 대상이면(shouldMaterialize) backing을 정명 경로로 promote(복사)한 뒤, 모든 경로에서 backing + overlay JSON 스크래치(.tsma/test)를 정리한다 — materialize 경로도 복사 후 삭제(W2 주 누적원), 실패/재시도 경로도 측정 후 즉시 삭제. backing은 측정 전용이고 재시도 Render는 디스크 정명 테스트만 읽으므로 보존할 필요가 없다.
package tsmagate

import (
	"github.com/park-jun-woo/reins/pkg/quest"
)

// finalizeBacking handles the overlay backing once measured (C2). When the result
// should be committed (shouldMaterialize) it first promotes (copies) the backing
// to the canonical path; then, on every path, it sweeps the backing and overlay
// JSON scratch under .tsma/test. The materialize path deletes after copying (the
// main W2 leak), and the failure/retry paths delete right after measurement. The
// backing is measurement-only and a retry Render reads only the canonical on-disk
// test, so nothing of value is lost.
func finalizeBacking(p funcPayload, it *quest.Item, m *measurement, backingRel string) {
	if shouldMaterialize(m, it) {
		promoteBacking(p, m, backingRel)
	}
	cleanupBacking(p.Root, backingRel)
}
