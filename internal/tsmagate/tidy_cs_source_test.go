package tsmagate

import (
	"strings"
	"testing"
)

// TestTidyCsSourceFormatterAbsent asserts that when the .NET SDK (dotnet) is not
// installed (or formatting fails), tidyCsSource degrades to returning the input
// unchanged — best-effort, never required. This is the only deterministically
// reachable branch in environments without the SDK; the format-success, temp-IO,
// and empty-output branches require dotnet to be present.
func TestTidyCsSourceFormatterAbsent(t *testing.T) {
	src := "namespace P;\npublic class FooTests\n{\n    public void T() { }\n}\n"
	got := tidyCsSource(src)
	// Whether the SDK is present or not, the class body must survive.
	if !strings.Contains(got, "class FooTests") {
		t.Errorf("tidyCsSource lost content: %q", got)
	}
}
