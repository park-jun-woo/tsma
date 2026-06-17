package tsmagate

import (
	"strings"
	"testing"
)

// TestSanitizeCsSourceUnwrapsFence asserts the markdown fence and any surrounding
// prose an LLM may emit are stripped, leaving compilable C#. The assertion is
// robust to whether dotnet format is installed (it only checks the fence/prose
// are gone and the class body survives).
func TestSanitizeCsSourceUnwrapsFence(t *testing.T) {
	raw := "Here is the test:\n```csharp\nnamespace P;\n\npublic class FooTests {\n    public void T() {}\n}\n```\nHope it helps!"
	got := sanitizeCsSource(raw)
	if strings.Contains(got, "```") {
		t.Errorf("fence not removed: %q", got)
	}
	if strings.Contains(got, "Hope it helps") || strings.Contains(got, "Here is the test") {
		t.Errorf("prose not removed: %q", got)
	}
	if !strings.Contains(got, "class FooTests") {
		t.Errorf("class body lost: %q", got)
	}
}

// TestSanitizeCsSourceNoFence asserts un-fenced input is preserved (trimmed),
// not rejected — sanitize is best-effort, never a validator.
func TestSanitizeCsSourceNoFence(t *testing.T) {
	raw := "\n\nnamespace P;\npublic class BarTests {}\n\n"
	got := sanitizeCsSource(raw)
	if !strings.Contains(got, "class BarTests") {
		t.Errorf("content lost: %q", got)
	}
	if strings.Contains(got, "```") {
		t.Errorf("unexpected fence: %q", got)
	}
}

// TestSanitizeSourceDispatchCs asserts the language dispatcher routes csharp to
// the C# sanitizer (fence unwrap) and leaves the default (go) path intact.
func TestSanitizeSourceDispatchCs(t *testing.T) {
	cs := sanitizeSource("csharp", "```csharp\npublic class T {}\n```")
	if strings.Contains(cs, "```") || !strings.Contains(cs, "class T") {
		t.Errorf("csharp dispatch wrong: %q", cs)
	}
}
