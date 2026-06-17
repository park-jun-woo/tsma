//ff:func feature=gate type=helper control=sequence lang=rust
//ff:what sanitizeRsSource: the Rust analogue of sanitizeCsSource — strips the markdown fences (```rust / ```) and out-of-fence prose an LLM may wrap a generated `#[cfg(test)] mod tests { ... }` block in, leaving pure Rust, then runs tidyRsSource (rustfmt if available, else identity). Unwrap is common (shared shape with Go/C#); formatting is best-effort. Never rejects — a still-broken block is caught downstream by the truncation check (parseRsTestFuncs) and tests-must-pass.
package tsmagate

import "strings"

// sanitizeRsSource unwraps Markdown code fences around a generated Rust test and
// passes the result through tidyRsSource (rustfmt best-effort). When no fence is
// present it trims surrounding whitespace. Best-effort, not a validator.
func sanitizeRsSource(raw string) string {
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
		return tidyRsSource(strings.TrimSpace(rest) + "\n")
	}
	return tidyRsSource(strings.TrimSpace(s) + "\n")
}
