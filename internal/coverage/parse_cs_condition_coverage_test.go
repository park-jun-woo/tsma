package coverage

import "testing"

func TestParseCsConditionCoverage(t *testing.T) {
	cases := []struct {
		in             string
		covered, total int
	}{
		{"50% (1/2)", 1, 2},
		{"100% (4/4)", 4, 4},
		{"0% (0/2)", 0, 2},
		{"", 0, 0},
		{"garbage", 0, 0},
		{"(x/y)", 0, 0},
		{"100% (noslash)", 0, 0}, // parens present but no '/' -> slash<0 branch
	}
	for _, tc := range cases {
		c, tot := parseCsConditionCoverage(tc.in)
		if c != tc.covered || tot != tc.total {
			t.Errorf("parseCsConditionCoverage(%q) = (%d,%d), want (%d,%d)",
				tc.in, c, tot, tc.covered, tc.total)
		}
	}
}
