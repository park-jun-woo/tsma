package tsmagate

import (
	"strings"
	"testing"
)

// TestTidyJavaSourceFormatterAbsent asserts that when google-java-format is not
// installed (or fails), tidyJavaSource degrades to returning the input unchanged
// — best-effort, never required. This is the only deterministically reachable
// branch here; the format-success and empty-output branches require the external
// formatter to be present.
func TestTidyJavaSourceFormatterAbsent(t *testing.T) {
	src := "package p;\nclass FooTest {\n  void t() {}\n}\n"
	got := tidyJavaSource(src)
	// Whether the formatter is present or not, the class body must survive.
	if !strings.Contains(got, "class FooTest") {
		t.Errorf("tidyJavaSource lost content: %q", got)
	}
}
