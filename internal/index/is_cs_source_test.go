package index

import "testing"

func TestIsCsSource(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"src/Foo.cs", true},
		{"Foo.cs", true},
		{"Foo.java", false},
		{"FooTest.cs", false},
		{"FooTests.cs", false},
		{"App.Tests/FooTests.cs", false},
		{"App.Tests/Helper.cs", false},
		{"App.Test/Helper.cs", false},
		{"src/App.csproj", false},
	}
	for _, tc := range cases {
		if got := isCsSource(tc.path); got != tc.want {
			t.Errorf("isCsSource(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}
