//ff:func feature=gate type=helper control=sequence lang=java
//ff:what sanitizeJavaSource: the Java analogue of sanitizeTSSource — strips the markdown fences (```java / ```) and out-of-fence prose an LLM may wrap a generated JUnit test in, leaving pure Java, then runs tidyJavaSource (google-java-format if available, else identity). Unwrap is common (shared shape with Go/TS); formatting is best-effort. Never rejects — a still-broken file is caught downstream by tests-must-pass.
package tsmagate

import "strings"

// sanitizeJavaSource unwraps Markdown code fences around a generated Java test
// and passes the result through tidyJavaSource (google-java-format best-effort).
// When no fence is present it trims surrounding whitespace. Best-effort, not a
// validator.
func sanitizeJavaSource(raw string) string {
	s := strings.ReplaceAll(raw, "\r\n", "\n")
	if i := strings.Index(s, "```"); i >= 0 {
		rest := s[i+3:]
		// Drop the optional language tag on the opening fence line.
		if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
			rest = rest[nl+1:]
		}
		if j := strings.Index(rest, "```"); j >= 0 {
			rest = rest[:j]
		}
		return tidyJavaSource(strings.TrimSpace(rest) + "\n")
	}
	return tidyJavaSource(strings.TrimSpace(s) + "\n")
}
