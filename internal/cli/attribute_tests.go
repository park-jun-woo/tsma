//ff:func feature=cli type=helper control=sequence
//ff:what Populates each function's TestFiles/TestFile via content-aware (Go) or fallback matching
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// attributeTests fills in TestFiles (full set) and TestFile (representative =
// TestFiles[0]) for every function. For Go it uses the content-aware batch
// matcher, which builds each package's index exactly once and reuses it for all
// functions in that package, so a single analyze parses each package only once.
// For every other language it falls back to the per-function file-name matcher,
// preserving existing single-file behavior.
func attributeTests(projectRoot, lang string, functions []model.Function) {
	if lang == "go" {
		attributeGoTests(projectRoot, functions)
		return
	}
	attributeFallbackTests(projectRoot, lang, functions)
}
