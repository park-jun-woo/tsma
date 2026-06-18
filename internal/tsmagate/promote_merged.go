//ff:func feature=gate type=helper control=sequence level=error lang=go
//ff:what promoteMerged: BUG-002 누적 promote 공통 헬퍼. 정명 경로(canonical)의 기존 내용을 읽어(없으면 빈) 새 src를 mergeCanonical(qn, lang)로 함수별 블록 누적한 뒤 writeTestFile로 그대로 쓴다(머지 후 재포맷 없음 — 마커 생존). 네이티브 promote 3종(Go/TS/Py)과 제네릭 prepare.go 제네릭 loop write가 공유한다. 읽기 실패(파일 부재 제외)·쓰기 실패는 error로 전파해 호출처가 TestFailed로 드러낸다(무음 금지). writeTestFile은 순수 덮어쓰기로 유지되고, 누적은 이 호출처-측 read→merge→write로만 일어난다(backing 측정용 write는 비누적).

package tsmagate

import (
	"os"
	"path/filepath"
)

// promoteMerged reads the existing canonical file (treating a missing file as
// empty), accumulates src into it as the qn function's marker block via
// mergeCanonical, and writes the merged result verbatim with writeTestFile (no
// re-format, so markers survive). This is the single read→merge→write path shared
// by the native promotes (Go/TS/Py) and the generic prepare.go loop write. A read
// error other than not-exist, or a write error, is returned so the caller can
// surface it as TestFailed (never silent).
func promoteMerged(root, canonical, src, qn, lang string) error {
	existing, err := os.ReadFile(filepath.Join(root, canonical))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	merged := mergeCanonical(string(existing), src, qn, lang)
	return writeTestFile(root, canonical, merged)
}
