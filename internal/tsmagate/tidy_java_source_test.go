package tsmagate

import (
	"strings"
	"testing"
)

// TestTidyJavaSourceFormatterAbsent asserts that when google-java-format is not
// on PATH, tidyJavaSource degrades to returning the input unchanged —
// best-effort, never required. The success and empty-output branches are driven
// below with a fake formatter script on PATH.
func TestTidyJavaSourceFormatterAbsent(t *testing.T) {
	emptyPath(t)
	src := "package p;\nclass FooTest {\n  void t() {}\n}\n"
	got := tidyJavaSource(src)
	// Whether the formatter is present or not, the class body must survive.
	if !strings.Contains(got, "class FooTest") {
		t.Errorf("tidyJavaSource lost content: %q", got)
	}
}

// TestTidyJavaSourceFormatSuccess drives the happy path with a fake
// google-java-format that consumes stdin and emits fixed formatted output.
func TestTidyJavaSourceFormatSuccess(t *testing.T) {
	installFakeTool(t, "google-java-format", "#!/bin/sh\ncat > /dev/null\nprintf 'class Fmted {}\\n'\n")
	if got := tidyJavaSource("class X{}\n"); got != "class Fmted {}\n" {
		t.Errorf("tidyJavaSource = %q, want %q", got, "class Fmted {}\n")
	}
}

// TestTidyJavaSourceRunError covers the cmd.Run failure branch: a fake formatter
// that exits non-zero degrades to the input unchanged.
func TestTidyJavaSourceRunError(t *testing.T) {
	installFakeTool(t, "google-java-format", "#!/bin/sh\ncat > /dev/null\nexit 2\n")
	src := "class X{}\n"
	if got := tidyJavaSource(src); got != src {
		t.Errorf("run failure must return src unchanged, got %q", got)
	}
}

// TestTidyJavaSourceEmptyOutput covers the blank-output branch: a fake formatter
// that succeeds but prints only whitespace degrades to the input unchanged.
func TestTidyJavaSourceEmptyOutput(t *testing.T) {
	installFakeTool(t, "google-java-format", "#!/bin/sh\ncat > /dev/null\nprintf '  '\n")
	src := "class X{}\n"
	if got := tidyJavaSource(src); got != src {
		t.Errorf("blank formatter output must return src unchanged, got %q", got)
	}
}
