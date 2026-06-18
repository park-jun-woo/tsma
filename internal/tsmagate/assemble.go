//ff:func feature=gate type=helper control=selection
//ff:what assemble: 누적 엔진(mergeCanonical)의 최종 조립 어댑터. 합쳐진 헤더 라인들과 본문을 하나의 정명 파일 텍스트로 잇는다 — 둘 다 비지 않으면 빈 줄 하나로 구분하고, 한쪽이 비면 다른 쪽만(가짜 빈 줄 없이) 돌려준다. 끝은 ensureTrailingNewline으로 정확히 개행 하나로 마감한다.

package tsmagate

import "strings"

// assemble joins a header block and a body into the final file text, separating
// them with a single blank line when both are non-empty.
func assemble(header []string, body string) string {
	h := strings.TrimRight(strings.Join(header, "\n"), "\n")
	b := strings.TrimLeft(body, "\n")
	switch {
	case h == "":
		return ensureTrailingNewline(b)
	case b == "":
		return ensureTrailingNewline(h)
	default:
		return ensureTrailingNewline(h + "\n\n" + b)
	}
}
