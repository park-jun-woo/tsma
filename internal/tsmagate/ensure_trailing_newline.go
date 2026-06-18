//ff:func feature=gate type=helper control=sequence
//ff:what ensureTrailingNewline: 누적 엔진 조립의 말미 정규화 헬퍼. 문자열 끝의 개행들을 모두 떼고 정확히 개행 하나를 붙여, 정명 파일이 항상 단일 개행으로 끝나게 한다(머지 누적 시 끝 개행이 불어나는 것 방지).

package tsmagate

import "strings"

// ensureTrailingNewline guarantees the file ends with exactly one newline.
func ensureTrailingNewline(s string) string {
	return strings.TrimRight(s, "\n") + "\n"
}
