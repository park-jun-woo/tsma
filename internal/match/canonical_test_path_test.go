//ff:func feature=match type=test
//ff:what CanonicalTestPath 단위테스트: Go 소스→<base>_test.go(같은 디렉터리) 도출, 비-Go·비-.go는 빈문자열, 그리고 GoMatcher.Match가 디스크 존재 시 같은 경로를 돌려주는 일관성을 검증한다.

package match

import (
	"path/filepath"
	"testing"
)

func TestCanonicalTestPath(t *testing.T) {
	cases := []struct {
		lang, src, want string
	}{
		{"go", filepath.Join("pkg", "foo.go"), filepath.Join("pkg", "foo_test.go")},
		{"go", "main.go", "main_test.go"},
		{"go", "notgo.txt", ""}, // not a .go source
		{"python", filepath.Join("a", "b.py"), ""}, // non-Go lang in Phase 002
	}
	for _, c := range cases {
		if got := CanonicalTestPath(c.lang, c.src); got != c.want {
			t.Errorf("CanonicalTestPath(%q,%q) = %q, want %q", c.lang, c.src, got, c.want)
		}
	}
}
