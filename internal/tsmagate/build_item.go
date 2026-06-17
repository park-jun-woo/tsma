//ff:func feature=gate type=helper control=sequence
//ff:what buildItem: Seed에서 분리한 함수→quest.Item 변환. reins가 소유하는 런타임 필드(Status/Attempt/CoveragePct/TestMtime/FailOutput)를 제거해 payload를 순수 위치 데이터로 만들고, Item.Key=QualifiedName·State=TODO로 만들어 lang/root/함수 스냅샷을 payload에 싣는다. Seed의 range body가 10줄을 넘지 않도록 추출(Q4).

package tsmagate

import (
	"fmt"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

// buildItem turns one indexed function into a seeded quest.Item. It strips the
// runtime fields reins now owns so the payload is pure location data
// (State/Tries are the ratchet's, not the payload's) and snapshots lang/root/fn
// into the item payload. Extracted from Seed so its range body stays short.
func buildItem(lang, root string, fn model.Function) (*quest.Item, error) {
	fn.Status = ""
	fn.Attempt = 0
	fn.CoveragePct = 0
	fn.TestMtime = ""
	fn.FailOutput = ""

	it := &quest.Item{Key: fn.QualifiedName, State: quest.TODO}
	if err := it.SetPayload(funcPayload{Lang: lang, Root: root, Fn: fn}); err != nil {
		return nil, fmt.Errorf("snapshot payload for %s: %w", fn.QualifiedName, err)
	}
	return it, nil
}
