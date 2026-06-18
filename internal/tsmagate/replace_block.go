//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what replaceOrAddBlock: 누적 엔진의 마커 블록 교체/추가 어댑터(D1). 기존 본문(헤더 제외)에서 `<line> tsma:begin fn=<QN>`…`<line> tsma:end fn=<QN>` 마커로 감싼 같은 함수 블록을 찾으면 그 구간만 새 블록으로 교체(같은 함수 재시도=그 함수 몫만 갱신), 없으면 기존 본문 뒤에 새 블록을 추가(다른 함수=누적)한다. 마커가 유실돼 같은 함수 블록을 못 찾으면 보수적으로 추가만 한다(소실<중복 폴백, Phase007 §4). 블록 식별은 라인-span이 아니라 마커 문자열로만 한다(model.Function은 테스트 span을 갖지 않음).

package tsmagate

import "strings"

const (
	markerBegin = "tsma:begin fn="
	markerEnd   = "tsma:end fn="
)

// replaceOrAddBlock returns the existing body with qn's marker block replaced by
// newBlock, or — when no block for qn exists (new function, or marker lost) —
// existingBody followed by newBlock. Identification is purely by the begin/end
// marker text for qn; on a lost marker it conservatively appends (duplication is
// preferred over silent loss, Phase007 §4).
func replaceOrAddBlock(existingBody, newBlock, qn string, tok commentTokens) string {
	begin := tok.line + " " + markerBegin + qn
	end := tok.line + " " + markerEnd + qn
	lines := strings.Split(existingBody, "\n")

	startIdx, endIdx := -1, -1
	for i, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if startIdx == -1 && trimmed == begin {
			startIdx = i
			continue
		}
		if startIdx != -1 && trimmed == end {
			endIdx = i
			break
		}
	}

	// No intact block for qn: append (new function or lost marker fallback).
	if startIdx == -1 || endIdx == -1 {
		base := strings.TrimRight(existingBody, "\n")
		if base == "" {
			return newBlock
		}
		return base + "\n\n" + newBlock
	}

	// Replace the [startIdx, endIdx] block (inclusive) with newBlock's lines.
	var out []string
	out = append(out, lines[:startIdx]...)
	out = append(out, strings.Split(strings.TrimRight(newBlock, "\n"), "\n")...)
	out = append(out, lines[endIdx+1:]...)
	return strings.Join(out, "\n")
}
