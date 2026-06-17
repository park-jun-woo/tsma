//ff:func feature=match type=helper control=iteration dimension=1 lang=typescript
//ff:what canonicalTSTestPath: the TypeScript arm of CanonicalTestPath (jest/vitest) — foo.ts → foo.test.ts in the same directory, original extension (.ts/.tsx/.js/.jsx) preserved; non-TS files return "". Same-dir matches TSMatcher.Match's 1st search dir so write and read agree.
package match

import (
	"path/filepath"
	"strings"
)

// canonicalTSTestPath inserts ".test" before the TS/JS extension of base, in the
// source's directory, or returns "" when base is not a TS/JS source.
func canonicalTSTestPath(sourceFile, base string) string {
	for _, ext := range []string{".tsx", ".jsx", ".ts", ".js"} {
		if strings.HasSuffix(base, ext) {
			testBase := strings.TrimSuffix(base, ext) + ".test" + ext
			return filepath.Join(filepath.Dir(sourceFile), testBase)
		}
	}
	return ""
}
