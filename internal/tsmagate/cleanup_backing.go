//ff:func feature=gate type=helper control=sequence lang=go
//ff:what cleanupBacking: overlay 측정이 끝난 뒤 .tsma/test의 스크래치를 지운다(C2/W2) — backing(gen-<item>.go)과 runner/checker가 공유하던 overlay.json 둘 다. 정리 주체는 tsmagate 호출부다(runner.Run·coverage.Check가 같은 overlay.json을 매번 재기록해 runner 자체삭제는 경합). 다음 시도는 raw에서 새 backing을 만들고 재시도 Render는 디스크 정명 테스트만 읽으므로 측정 후 삭제는 무해. 삭제 실패는 무음(스크래치라 치명 아님).
package tsmagate

import (
	"os"
	"path/filepath"
)

// cleanupBacking removes the overlay scratch under .tsma/test once a measurement
// is done (C2/W2): both the backing file (gen-<item>.go) and the overlay.json that
// runner.Run and coverage.Check share. The caller (tsmagate) owns this cleanup
// rather than the runner, because both rewrite the same overlay.json every
// measurement, so a runner self-delete would race the two calls. Deleting after
// the measurement is harmless: the next try rebuilds a fresh backing from raw and
// a retry Render reads only the canonical on-disk test, never the backing. Remove
// errors are ignored — this is scratch, not source of truth.
func cleanupBacking(root, backingRel string) {
	_ = os.Remove(filepath.Join(root, backingRel))
	_ = os.Remove(filepath.Join(root, ".tsma", "test", "overlay.json"))
}
