package match

import "testing"

func TestIsJavaTestFile(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"FooTest.java", true},
		{"FooTests.java", true},
		{"Foo.java", false},   // .java but stem is neither Test nor Tests
		{"FooTest.txt", false}, // not a .java file
		{"Test.java", true},    // stem ends with Test
	}
	for _, c := range cases {
		if got := isJavaTestFile(c.name); got != c.want {
			t.Errorf("isJavaTestFile(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}
