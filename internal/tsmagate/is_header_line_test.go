//ff:func feature=gate type=test
//ff:what isHeaderLine 단위테스트: 헤더 키워드 매치(import/package), 비헤더 라인,
// 접두사 집합이 빈 경우를 덮는다.
package tsmagate

import "testing"

func TestIsHeaderLine(t *testing.T) {
	tok := commentTokensFor("go")
	for _, line := range []string{`import "testing"`, "package tsmagate", "import\t\"x\""} {
		if !isHeaderLine(line, tok) {
			t.Errorf("isHeaderLine(%q) = false, want true", line)
		}
	}
	if isHeaderLine("func TestX(t *testing.T) {}", tok) {
		t.Error("a body line must not be a header line")
	}
	if isHeaderLine("import stuff", commentTokens{}) {
		t.Error("no header prefixes -> never a header line")
	}
}
