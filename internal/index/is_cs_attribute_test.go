package index

import "testing"

func TestIsCsAttribute(t *testing.T) {
	for _, line := range []string{"[Fact]", "[TestMethod]", `[Obsolete("x")]`} {
		if !isCsAttribute(line) {
			t.Errorf("isCsAttribute(%q) = false, want true", line)
		}
	}
	for _, line := range []string{"public int Add() {", "namespace Foo {", "// comment"} {
		if isCsAttribute(line) {
			t.Errorf("isCsAttribute(%q) = true, want false", line)
		}
	}
}
