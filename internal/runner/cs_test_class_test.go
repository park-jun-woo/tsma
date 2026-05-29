package runner

import (
	"path/filepath"
	"testing"
)

func TestCsTestClass(t *testing.T) {
	cases := map[string]string{
		filepath.Join("App.Tests", "Services", "FooTests.cs"): "FooTests",
		filepath.Join("src", "BarTest.cs"):                    "BarTest",
		"BazTests.cs":                                         "BazTests",
	}
	for in, want := range cases {
		if got := csTestClass(in); got != want {
			t.Errorf("csTestClass(%q) = %q, want %q", in, got, want)
		}
	}
}
