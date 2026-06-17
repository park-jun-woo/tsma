//ff:func feature=gate type=helper control=iteration dimension=1 lang=go
//ff:what firstMalformedTestName: 추출한 테스트 함수명 중 go test가 조용히 무시하는 첫 잘못된 이름을 돌려준다(C1 보강). Go 규칙상 "Test" 다음 글자가 소문자면(예: TestpyIndent) 테스트로 인식 안 돼 0개 실행→exit0 무용통과가 된다. 그런 이름이 있으면 (이름, true), 전부 정상이면 ("", false). 측정 전 reject로 0% 가드보다 구체적인 피드백을 준다.
package tsmagate

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// firstMalformedTestName returns the first extracted test-func name that go test
// would silently ignore — one where the rune right after the "Test" prefix is a
// lowercase letter (e.g. `TestpyIndent`), so it never runs and a 0-test exit 0
// masquerades as a pass (C1). ok is false when every name is well-formed. The
// bare "Test" and any non-"Test" name (parseTestFuncs only returns "Test"-prefixed
// ones) are treated as fine here.
func firstMalformedTestName(funcs []string) (string, bool) {
	for _, name := range funcs {
		rest := strings.TrimPrefix(name, "Test")
		if rest == "" || rest == name {
			continue
		}
		r, _ := utf8.DecodeRuneInString(rest)
		if unicode.IsLower(r) {
			return name, true
		}
	}
	return "", false
}
