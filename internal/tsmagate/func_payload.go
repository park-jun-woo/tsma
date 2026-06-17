//ff:type feature=gate type=model
//ff:what funcPayload는 quest.Item.Payload에 담는 도메인 스냅샷이다. Render/Prepare가 한 함수를 결정적으로 재탐색·재매칭·재측정하는 데 필요한 lang/root/함수위치를 운반한다. Seed가 reins 세션에 stash할 곳이 없어 lang·root를 모든 item에 중복 저장해 자기완결적으로 만든다. Status(→Item.State)·Attempt(→Item.Tries)는 reins가 소유하므로 seed 시점에 제거한다(래칫이 진행상태의 단일 진실).

package tsmagate

import "github.com/park-jun-woo/tsma/internal/model"

// funcPayload is the domain snapshot stored in quest.Item.Payload. It carries
// everything Render/Prepare need to re-locate, re-match, and re-measure a
// function deterministically. lang and root are duplicated into every item
// because Seed has no session to stash them in (reins seeds items, not Meta),
// keeping each item self-contained. Fn deliberately omits the runtime fields
// reins now owns: Status (→ Item.State) and Attempt (→ Item.Tries) are zeroed at
// seed time so the ratchet is the single source of truth for progress.
type funcPayload struct {
	Lang string         `json:"lang"`
	Root string         `json:"root"`
	Fn   model.Function `json:"fn"`
}
