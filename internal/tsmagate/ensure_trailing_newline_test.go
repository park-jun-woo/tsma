//ff:func feature=gate type=test
//ff:what ensureTrailingNewline 단위테스트: 개행 없음/하나/여럿/빈 문자열 모두 정확히
// 단일 개행으로 끝나는지 확인한다.
package tsmagate

import "testing"

func TestEnsureTrailingNewline(t *testing.T) {
	cases := map[string]string{
		"a":       "a\n",
		"a\n":     "a\n",
		"a\n\n\n": "a\n",
		"":        "\n",
	}
	for in, want := range cases {
		if got := ensureTrailingNewline(in); got != want {
			t.Errorf("ensureTrailingNewline(%q) = %q, want %q", in, got, want)
		}
	}
}
