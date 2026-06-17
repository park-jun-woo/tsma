//ff:func feature=match type=helper control=sequence lang=typescript
//ff:what tsFilenameFallback: attributes a TS/JS function to its conventional sibling test file via TSMatcher.Match when content-aware matching found nothing — the last-resort fallback the plan requires be preserved (not removed). TestFuncs nil = run the whole file.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// tsFilenameFallback returns the conventional <name>.test.ts (same dir or
// __tests__/) for fn via TSMatcher.Match, or found=false when none exists.
func tsFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	testRel, ok := (&TSMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
