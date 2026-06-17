//ff:func feature=gate type=helper control=sequence lang=csharp
//ff:what sanitizeCsSource: the C# analogue of sanitizeJavaSource — strips the markdown fences (```csharp / ```) and out-of-fence prose an LLM may wrap a generated xUnit test in, leaving pure C#, then runs tidyCsSource (dotnet format if available, else identity). Unwrap is common (shared shape with Go/Java); formatting is best-effort. Never rejects — a still-broken file is caught downstream by tests-must-pass.
package tsmagate

import "strings"

// sanitizeCsSource unwraps Markdown code fences around a generated C# test and
// passes the result through tidyCsSource (dotnet format best-effort). When no
// fence is present it trims surrounding whitespace. Best-effort, not a validator.
func sanitizeCsSource(raw string) string {
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
		return tidyCsSource(strings.TrimSpace(rest) + "\n")
	}
	return tidyCsSource(strings.TrimSpace(s) + "\n")
}
