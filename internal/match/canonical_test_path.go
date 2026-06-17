//ff:func feature=match type=helper control=selection
//ff:what Derives the canonical test path for a source file — the single source of truth for the test-naming formula reused by GoMatcher.Match, FindMisnamedTest, and the loop's write path. Go: <base>_test.go. TypeScript (Phase005a): <base>.test.<ext> (jest/vitest). Python (Phase005b): test_<base>.py (pytest). Java (Phase005c): src/test mirror FooTest.java (JUnit, via javaTestDir SSOT). Non-handled languages return "".
package match

import (
	"path/filepath"
	"strings"
)

// CanonicalTestPath derives the canonical test-file path for a source file in the
// given language: the conventional sibling test file that, by naming convention,
// covers it. For Go that is "<base>_test.go" in the same directory as the source
// (e.g. "pkg/foo.go" → "pkg/foo_test.go"). The returned path is in whatever form
// sourceFile was given (project-root-relative in, project-root-relative out).
//
// It returns "" when no canonical path can be derived — a non-Go language, or a
// source file that does not end in ".go". This is the one place the Go
// source→test naming formula lives; GoMatcher.Match and FindMisnamedTest both
// build on it, and the loop write path (Phase 002) uses it to place a brand-new
// test file. Only Go is handled in Phase 002; other languages return "".
func CanonicalTestPath(lang, sourceFile string) string {
	base := filepath.Base(sourceFile)
	switch lang {
	case "go":
		if !strings.HasSuffix(base, ".go") {
			return ""
		}
		testBase := strings.TrimSuffix(base, ".go") + "_test.go"
		return filepath.Join(filepath.Dir(sourceFile), testBase)
	case "typescript":
		return canonicalTSTestPath(sourceFile, base)
	case "python":
		return canonicalPyTestPath(sourceFile, base)
	case "java":
		return canonicalJavaTestPath(sourceFile, base)
	default:
		return ""
	}
}
