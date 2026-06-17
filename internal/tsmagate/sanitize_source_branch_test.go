package tsmagate

import (
	"strings"
	"testing"
)

// TestSanitizeSourceAllArms covers every language arm of the dispatcher
// (typescript, python, java, and the default/Go fallback). Each arm unwraps the
// markdown fence; formatters (prettier/black/google-java-format) are optional, so
// the assertions only check the fence is gone and the body survives.
func TestSanitizeSourceAllArms(t *testing.T) {
	cases := []struct {
		lang string
		raw  string
		body string
	}{
		{"typescript", "```ts\nexport const x = 1;\n```", "export const x = 1"},
		{"python", "```python\ndef test_x():\n    pass\n```", "def test_x"},
		{"java", "```java\nclass FooTest {}\n```", "class FooTest"},
		{"go", "```go\npackage p\n```", "package p"},
		{"rust", "```rust\nfn t() {}\n```", "fn t"}, // unknown lang -> default arm
	}
	for _, c := range cases {
		got := sanitizeSource(c.lang, c.raw)
		if strings.Contains(got, "```") {
			t.Errorf("%s: fence not removed: %q", c.lang, got)
		}
		if !strings.Contains(got, c.body) {
			t.Errorf("%s: body %q lost in %q", c.lang, c.body, got)
		}
	}
}
