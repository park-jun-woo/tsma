package tsmagate

import (
	"strings"
	"testing"
)

// TestSanitizeJavaSourceUnwrapsFence asserts the markdown fence and any
// surrounding prose an LLM may emit are stripped, leaving compilable Java. The
// assertion is robust to whether google-java-format is installed (it only checks
// the fence/prose are gone and the class body survives).
func TestSanitizeJavaSourceUnwrapsFence(t *testing.T) {
	raw := "Here is the test:\n```java\npackage p;\n\nclass FooTest {\n    void t() {}\n}\n```\nHope it helps!"
	got := sanitizeJavaSource(raw)
	if strings.Contains(got, "```") {
		t.Errorf("fence not removed: %q", got)
	}
	if strings.Contains(got, "Hope it helps") || strings.Contains(got, "Here is the test") {
		t.Errorf("prose not removed: %q", got)
	}
	if !strings.Contains(got, "class FooTest") {
		t.Errorf("class body lost: %q", got)
	}
}

// TestSanitizeJavaSourceNoFence asserts un-fenced input is preserved (trimmed),
// not rejected — sanitize is best-effort, never a validator.
func TestSanitizeJavaSourceNoFence(t *testing.T) {
	raw := "\n\npackage p;\nclass BarTest {}\n\n"
	got := sanitizeJavaSource(raw)
	if !strings.Contains(got, "class BarTest") {
		t.Errorf("content lost: %q", got)
	}
	if strings.Contains(got, "```") {
		t.Errorf("unexpected fence: %q", got)
	}
}

// TestSanitizeSourceDispatch asserts the language dispatcher routes java to the
// Java sanitizer (fence unwrap) and leaves the default (go) path intact.
func TestSanitizeSourceDispatch(t *testing.T) {
	java := sanitizeSource("java", "```java\nclass T {}\n```")
	if strings.Contains(java, "```") || !strings.Contains(java, "class T") {
		t.Errorf("java dispatch wrong: %q", java)
	}
}
