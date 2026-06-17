//ff:func feature=gate type=helper control=sequence lang=python
//ff:what sanitizePySource: the Python analogue of sanitizeGoSource/sanitizeTSSource — strips the markdown fences (```python / ```py / ```) and out-of-fence prose an LLM may wrap a generated test in, leaving pure Python, then runs tidyPySource (black/isort if available, else identity). Unwrap is common; formatting is best-effort. Never rejects — a still-broken file is caught downstream by tests-must-pass.
package tsmagate

import "strings"

// sanitizePySource unwraps Markdown code fences around a generated Python test
// and passes the result through tidyPySource (black/isort best-effort). When no
// fence is present it trims surrounding whitespace. Best-effort, not a validator.
func sanitizePySource(raw string) string {
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
		return tidyPySource(strings.TrimSpace(rest) + "\n")
	}
	return tidyPySource(strings.TrimSpace(s) + "\n")
}
