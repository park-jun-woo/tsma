//ff:func feature=gate type=test
//ff:what firstMalformedTestName 단위테스트: 소문자-후속(TestpyIndent)은 잡고, 정상명(TestFoo)·대문자/숫자/언더스코어 후속·바로 "Test"는 통과, 첫 위반을 순서대로 돌려줌을 확인한다.
package tsmagate

import "testing"

func TestFirstMalformedTestName(t *testing.T) {
	cases := []struct {
		name    string
		funcs   []string
		wantBad string
		wantOK  bool
	}{
		{"lowercase-after-test", []string{"TestpyIndent"}, "TestpyIndent", true},
		{"well-formed", []string{"TestFoo", "TestBar"}, "", false},
		{"bare-test-ok", []string{"Test"}, "", false},
		{"underscore-ok", []string{"Test_helper"}, "", false},
		{"digit-ok", []string{"Test1"}, "", false},
		{"first-violation-wins", []string{"TestGood", "TestbadOne", "TestbadTwo"}, "TestbadOne", true},
		{"empty", nil, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bad, ok := firstMalformedTestName(c.funcs)
			if bad != c.wantBad || ok != c.wantOK {
				t.Fatalf("firstMalformedTestName(%v) = (%q,%v), want (%q,%v)", c.funcs, bad, ok, c.wantBad, c.wantOK)
			}
		})
	}
}
