//ff:func feature=gate type=helper control=selection
//ff:what sanitizeSource: dispatches generated-test cleanup to the language's sanitizer — Go via gofmt (sanitizeGoSource), TypeScript via prettier (sanitizeTSSource), Python via black/isort (sanitizePySource), Java via google-java-format (sanitizeJavaSource), C# via dotnet format (sanitizeCsSource). Every arm unwraps markdown fences first (common) and formats best-effort. Unknown languages fall back to the Go sanitizer's fence-unwrap, preserving the prior generic disk-write behavior. Used by Prepare's generic loop write path so each language formats its own scratch test.
package tsmagate

// sanitizeSource cleans a generated test for the given language (fence unwrap +
// best-effort format). It falls back to sanitizeGoSource (fence unwrap) for
// languages without a dedicated sanitizer, matching prior behavior.
func sanitizeSource(lang, raw string) string {
	switch lang {
	case "typescript":
		return sanitizeTSSource(raw)
	case "python":
		return sanitizePySource(raw)
	case "java":
		return sanitizeJavaSource(raw)
	case "csharp":
		return sanitizeCsSource(raw)
	case "rust":
		return sanitizeRsSource(raw)
	default:
		return sanitizeGoSource(raw)
	}
}
