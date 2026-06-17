//ff:func feature=gate type=helper control=sequence lang=go
//ff:what finalizeBacking: overlay 측정 종결 처리. 확정 대상이면(shouldMaterialize) backing을 정명 경로로 promote하고, 아니면 마지막 시도 TestFailed에 한해 backing(.tsma/test)을 폐기한다. 그 외(재시도 예정)는 backing을 보존해 다음 턴 컨텍스트로 남길 수 있게 둔다.
package tsmagate

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// finalizeBacking decides what happens to the overlay backing once measured. When
// the result should be committed (shouldMaterialize) it promotes the backing to
// the canonical path; otherwise it discards the backing only on the final-try
// failure (a broken last artifact). Any other case (a retry is coming) keeps the
// backing under .tsma/test so it remains available.
func finalizeBacking(p funcPayload, it *quest.Item, m *measurement, backingRel string) {
	if shouldMaterialize(m, it) {
		promoteBacking(p, m, backingRel)
		return
	}
	if m.TestFailed && it.Tries == quest.MaxTries-1 {
		_ = os.Remove(filepath.Join(p.Root, backingRel))
	}
}
