//ff:func feature=match type=implementation control=iteration dimension=1 lang=go
//ff:what Attributes a Go function to its tests by content-aware package index lookup
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc builds the content-aware test index for the function's package once
// and looks the function up by its bare name. It returns the deduplicated set of
// test files and test functions that reference the function. When no test in the
// package references it (or the package cannot be indexed), it returns found
// false. It deliberately does not fall back to file-name matching: pure
// content-aware attribution is the contract here, and any fallback is handled by
// the caller in a later phase.
func (m *GoFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	pkgDir := filepath.Dir(fn.File)
	idx, err := BuildPkgTestIndex(projectRoot, pkgDir)
	if err != nil {
		return TestMatch{}, false
	}
	refs, ok := MatchFuncByName(idx, fn)
	if !ok {
		return TestMatch{}, false
	}
	return refsToTestMatch(refs)
}
