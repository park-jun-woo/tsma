//ff:func feature=gate type=helper control=sequence
//ff:what firstLines: s의 최대 n줄만(트림해) 돌려줘 시끄러운 컴파일러 덤프가 Fact를 범람시키지 않게 한다. 초과 시 n줄 뒤에 "…"를 덧붙인다.

package tsmagate

import "strings"

// firstLines returns at most n lines of s, trimmed, so a noisy compiler dump does
// not flood the Fact.
func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append(lines[:n], "…")
	}
	return strings.Join(lines, "\n")
}
