//ff:func feature=match type=helper control=sequence lang=go
//ff:what Falls back to file-name matching for a Go function via its <base>_test.go
package match

import "github.com/park-jun-woo/tsma/internal/model"

// goFilenameFallback attributes a Go function to its conventional
// <base>_test.go file when content-aware matching found nothing. It uses
// GoMatcher.Match (same match package, so no import cycle) to confirm the
// conventional test file exists on disk; when it does it returns a single-file
// TestMatch with TestFuncs left nil (the runner resolves the test functions
// from the file, identical to the non-Go fallback). When no conventional test
// file exists it reports found false. This is the second-pass fallback that
// keeps indirect/dispatch-style tests (e.g. cobra runCmd(t,"agent")) from
// producing false TODOs while preserving content-aware precision for functions
// that are directly referenced.
func goFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	testRel, ok := (&GoMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
