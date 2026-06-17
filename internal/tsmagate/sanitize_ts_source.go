//ff:func feature=gate type=helper control=sequence lang=typescript
//ff:what sanitizeTSSource: the TS analogue of sanitizeGoSource — strips the markdown fences (```ts / ```typescript / ```) and out-of-fence prose an LLM may wrap a generated test in, leaving pure TS, then runs tidyTSSource (prettier if available, else identity). Unwrap is common (shared with Go's shape); formatting is best-effort. Never rejects — a still-broken file is caught downstream by tests-must-pass.
package tsmagate

import "strings"

// sanitizeTSSource unwraps Markdown code fences around a generated TS test and
// passes the result through tidyTSSource (prettier best-effort). When no fence
// is present it trims surrounding whitespace. Best-effort, not a validator.
func sanitizeTSSource(raw string) string {
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
		return tidyTSSource(strings.TrimSpace(rest) + "\n")
	}
	return tidyTSSource(strings.TrimSpace(s) + "\n")
}
