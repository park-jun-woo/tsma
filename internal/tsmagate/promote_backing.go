//ff:func feature=gate type=helper control=sequence level=error lang=go
//ff:what promoteBacking: 통과(또는 통과DONE 예정)한 overlay backing 파일을 정명 경로(match.CanonicalTestPath)로 확정 기록한다 — 종결 시에만 실파일이 소스 트리에 닿는다. BUG-002: 통째 덮어쓰기 대신 promoteMerged로 함수별(QualifiedName) 마커 블록 누적 — 같은 소스 파일 다함수가 서로 덮어쓰지 않는다. backing 읽기·정명 누적쓰기 실패는 m.TestFailed로 드러내(무음 금지) PASS가 디스크에 남지 못하면 FAIL로 되먹임한다.
package tsmagate

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// promoteBacking commits a passing (or about-to-DONE) overlay backing file to its
// canonical test path — the one moment a real file reaches the source tree.
// BUG-002: instead of a whole-file overwrite it accumulates the backing as the
// function's marker block via promoteMerged, so multiple functions in one source
// file no longer overwrite each other. A read of the backing or an accumulating
// write to the canonical path that fails is surfaced as m.TestFailed (never
// silent): a PASS that cannot persist to disk is fed back as a FAIL.
func promoteBacking(p funcPayload, m *measurement, backingRel string) {
	canonical := match.CanonicalTestPath(p.Lang, p.Fn.File)
	data, err := os.ReadFile(filepath.Join(p.Root, backingRel))
	if err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
		return
	}
	if err := promoteMerged(p.Root, canonical, string(data), p.Fn.QualifiedName, p.Lang); err != nil {
		m.TestFailed, m.FailOutput = true, err.Error()
	}
}
