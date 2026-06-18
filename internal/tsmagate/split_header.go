//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what splitHeader: 누적 엔진의 헤더/본문 분리 어댑터. 한 완전 테스트 src를 (헤더 라인들, 본문)으로 가른다 — 헤더는 파일 선두의 연속 구간 중 빈 줄·라인 주석·tok.headerPrefixes(import/package/use/using…)로 시작하는 라인까지이고, 첫 비-헤더 비-주석 라인부터가 본문이다. 본문은 다시 합쳐 fn= 마커 블록으로 감쌀 원자료가 된다. tsma:begin/end 마커 라인은 헤더로 오인하지 않도록 본문 경계로 취급한다(폴백 append 시 마커 중첩 방지). 라인을 순회하며 헤더 경계를 찾는 iteration.

package tsmagate

import "strings"

// splitHeader divides a complete generated test source into its header lines (the
// leading import/package/use region, plus blank lines and line comments that sit
// inside it) and the remaining body. The header is the longest leading run of
// lines that are blank, a line comment, or start (after trimming) with one of
// tok.headerPrefixes; the first non-header, non-comment, non-blank line begins the
// body. A tsma:begin/end marker line (a comment, detected by substring not
// prefix) always ends the header region so existing accumulated blocks are never
// re-absorbed as header.
func splitHeader(src string, tok commentTokens) (header []string, body string) {
	lines := strings.Split(src, "\n")
	i := 0
	for ; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			header = append(header, lines[i])
			continue
		}
		if strings.Contains(trimmed, markerBegin) || strings.Contains(trimmed, markerEnd) {
			break
		}
		if strings.HasPrefix(trimmed, tok.line) {
			header = append(header, lines[i])
			continue
		}
		if isHeaderLine(trimmed, tok) {
			header = append(header, lines[i])
			continue
		}
		break
	}
	body = strings.Join(lines[i:], "\n")
	return header, body
}
