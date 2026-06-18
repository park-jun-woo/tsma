//ff:func feature=gate type=helper control=sequence
//ff:what wrapBlock: 누적 엔진의 마커 블록 생성 헬퍼(D1). 한 테스트 본문을 그 언어의 라인-주석 토큰으로 `<line> tsma:begin fn=<QN>` … `<line> tsma:end fn=<QN>` 마커로 감싸 블록을 만든다. 이렇게 만든 블록은 나중에 같은 함수 promote가 정확히 이 구간을 찾아 교체할 수 있게 한다(마커가 유일한 함수↔테스트 귀속 수단). 본문 끝의 잉여 개행은 떼어 마커 사이 간격을 일정하게 둔다.

package tsmagate

import "strings"

// wrapBlock wraps a test body in the per-function markers for qn using the
// language's line-comment token. The returned block is bounded by a leading
// begin-marker line and a trailing end-marker line so a later promote of the same
// function can locate and replace exactly this region.
func wrapBlock(body, qn string, tok commentTokens) string {
	var b strings.Builder
	b.WriteString(tok.line + " " + markerBegin + qn + "\n")
	b.WriteString(strings.TrimRight(body, "\n"))
	b.WriteString("\n" + tok.line + " " + markerEnd + qn + "\n")
	return b.String()
}
