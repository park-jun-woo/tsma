//ff:func feature=gate type=helper control=sequence level=error
//ff:what writeLoopSubmission: loop 모드 전용 생성테스트 디스크 기록. §2-1/§2-2 순서대로 사전 매칭(match.NewFuncMatcher)→경로 도출(testTargetPath)→정제(sanitizeSource)→함수별 블록 누적 기록(promoteMerged)한다. 경로 도출·쓰기 실패는 error로 돌려 Prepare가 TestFailed로 드러낸다(무음 금지). 성공 후 Prepare의 재매칭→실행→측정 공통 경로에 합류한다.

package tsmagate

import "github.com/park-jun-woo/tsma/internal/match"

// writeLoopSubmission writes the LLM-generated test (raw) to disk for loop mode,
// following the §2-1/§2-2 order: derive the target path from a pre-write match →
// sanitize → write as this function's accumulated marker block (promoteMerged,
// BUG-002). A path-derivation or write failure is returned as an error so the
// caller (Prepare) can surface it as TestFailed — never silent. On success the
// caller falls through to the shared re-match → run → measure path.
func writeLoopSubmission(p funcPayload, raw []byte) error {
	pre, preFound := match.NewFuncMatcher(p.Lang).MatchFunc(p.Root, &p.Fn)
	path, err := testTargetPath(p, pre, preFound)
	if err != nil {
		return err
	}
	return promoteMerged(p.Root, path, sanitizeSource(p.Lang, string(raw)), p.Fn.QualifiedName, p.Lang)
}
