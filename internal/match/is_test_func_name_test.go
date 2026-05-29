package match

import "testing"

func TestIsTestFuncName(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"TestFoo", true},
		{"TestX", true},
		{"Test_", true},
		{"Test", false},  // bare Test is excluded
		{"Tes", false},
		{"", false},
		{"testFoo", false},     // lowercase, not an entrypoint
		{"BenchmarkFoo", false},
		{"helperTest", false},  // Test not a prefix
	}
	for _, c := range cases {
		if got := isTestFuncName(c.name); got != c.want {
			t.Errorf("isTestFuncName(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}
