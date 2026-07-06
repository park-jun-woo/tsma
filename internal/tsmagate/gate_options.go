//ff:func feature=gate type=helper control=sequence
//ff:what GateOptions: tsma 판정을 TANGEUL 게이트 문서(gate.md, go:embed)로 옵트인하는 cli.GateOptions를 만든다. Case는 기본값과 동일한 "제출 통과"를 명시해 감사 표면을 고정하고, Registry는 gateRegistry()가 바인딩한다. tsma 술어는 ground 미사용(전부 tier-0)이라 Provider/Resolver는 nil. main이 cli.Options{Gate}에 끼운다.

package tsmagate

import (
	_ "embed"

	"github.com/park-jun-woo/reins/pkg/cli"
)

//go:embed gate.md
var gateDoc []byte

// GateOptions opts tsma's judgment into the TANGEUL gate document. tsma's
// predicates have no needs, so Provider and Resolver stay nil (all tier-0).
func GateOptions() *cli.GateOptions {
	return &cli.GateOptions{
		Src:      gateDoc,
		Path:     "internal/tsmagate/gate.md",
		Case:     "제출 통과", // 기본값과 동일 — 명시로 감사 표면 고정
		Registry: gateRegistry(),
		// Provider: nil — ground 미사용
		// Resolver: nil — 네트워크 seam 불필요
	}
}
