//ff:func feature=gate type=helper control=sequence level=error lang=go
//ff:what promoteBacking: 통과(또는 통과DONE 예정)한 overlay backing 파일을 정명 경로(match.CanonicalTestPath)로 확정 기록한다 — 종결 시에만 실파일이 소스 트리에 닿는다. backing 읽기·정명 쓰기 실패는 m.TestFailed로 드러내(무음 금지) PASS가 디스크에 남지 못하면 FAIL로 되먹임한다.
package tsmagate

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// promoteBacking commits a passing (or about-to-DONE) overlay backing file to its
// canonical test path — the one moment a real file reaches the source tree. A
// read of the backing or a write to the canonical path that fails is surfaced as
// m.TestFailed (never silent): a PASS that cannot persist to disk is fed back as
// a FAIL rather than vanishing.
func promoteBacking(p funcPayload, m *measurement, backingRel string) {
	canonical := match.CanonicalTestPath(p.Lang, p.Fn.File)
	data, err := os.ReadFile(filepath.Join(p.Root, backingRel))
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return
	}
	if err := writeTestFile(p.Root, canonical, string(data)); err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
	}
}
