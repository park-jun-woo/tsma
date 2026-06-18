//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what isHeaderLine: splitHeader의 헤더 판정 보조. trim된 라인이 그 언어의 헤더 키워드(import/package/use/using/from…) 중 하나로 시작하는지를 tok.headerPrefixes를 순회하며 판정한다 — 헤더(import) 영역과 테스트 본문을 가르는 1차 라인 단위 규칙.

package tsmagate

import "strings"

// isHeaderLine reports whether a trimmed line begins with one of the language's
// header keywords (import/package/use/using/from…).
func isHeaderLine(trimmed string, tok commentTokens) bool {
	for _, p := range tok.headerPrefixes {
		if strings.HasPrefix(trimmed, p) {
			return true
		}
	}
	return false
}
