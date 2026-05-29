//ff:func feature=match type=helper control=sequence lang=go
//ff:what Attributes one Go function content-aware first, then file-name fallback
package match

import "github.com/park-jun-woo/tsma/internal/model"

// attributeFunc resolves the TestMatch for a single Go function using the
// prebuilt package index idx (which may be nil when the package could not be
// indexed). It applies the hybrid rule: content-aware attribution wins when a
// test references the function by name; otherwise it falls back to file-name
// matching of the conventional <base>_test.go via goFilenameFallback. Returning
// found false means neither path produced a match. This is the shared per-
// function step of the batch MatchFuncs path, mirroring GoFuncMatcher.MatchFunc.
func attributeFunc(projectRoot string, idx *PkgTestIndex, fn *model.Function) (TestMatch, bool) {
	if refs, ok := MatchFuncByName(idx, fn); ok {
		if tm, ok := refsToTestMatch(refs); ok {
			return tm, true
		}
	}
	return goFilenameFallback(projectRoot, fn)
}
