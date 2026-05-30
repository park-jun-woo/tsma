//ff:func feature=match type=implementation control=iteration dimension=1 lang=go
//ff:what Attributes a Go function to its tests by content-aware package index lookup
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc builds the content-aware test index and the source-receiver map for
// the function's package once and looks the function up by its bare name with
// receiver-aware filtering. It returns the deduplicated set of test files and
// test functions that reference the function (content-aware precision: 1:N
// multi-file attribution is preserved; same-named methods on different
// receivers are kept apart). When content-aware
// matching finds nothing — the package cannot be indexed, no test references the
// function, or the refs resolve to no files — it falls back to file-name
// matching via goFilenameFallback, attributing the conventional <base>_test.go
// when it exists on disk. This hybrid keeps indirect/dispatch-style tests from
// producing false TODOs while never overriding a content-aware match. It mirrors
// the batch MatchFuncs path so detectTestChange (single-func) and analyze
// (batch) re-match identically.
func (m *GoFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	pkgDir := filepath.Dir(fn.File)
	idx, err := BuildPkgTestIndex(projectRoot, pkgDir)
	if err != nil {
		return goFilenameFallback(projectRoot, fn)
	}
	// Build the source-receiver map for the same package directory once, exactly
	// as the batch path does, so single-func re-matching resolves same-name-
	// multiple identically (a read error leaves it nil -> non-regressing).
	srcReceivers, _ := BuildPkgSourceReceivers(projectRoot, pkgDir)
	refs, ok := MatchFuncByName(idx, srcReceivers, fn)
	if !ok {
		return goFilenameFallback(projectRoot, fn)
	}
	if tm, ok := refsToTestMatch(refs); ok {
		return tm, true
	}
	return goFilenameFallback(projectRoot, fn)
}
