//ff:func feature=match type=implementation control=sequence
//ff:what Wraps the legacy Matcher.Match result into a single-file TestMatch
package match

import "github.com/park-jun-woo/tsma/internal/model"

// MatchFunc delegates to the wrapped file-name based Matcher using the
// function's source file and wraps the single matched test file into a
// TestMatch. TestFuncs is left nil so that the match package does not depend on
// the runner package (which would create an import cycle); the runner resolves
// the test functions from the file when TestFuncs is nil. Behavior is identical
// to the legacy Matcher.Match path.
func (m *fallbackFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil || m.inner == nil {
		return TestMatch{}, false
	}
	testFile, ok := m.inner.Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testFile}, TestFuncs: nil}, true
}
