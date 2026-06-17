//ff:func feature=match type=helper control=sequence lang=python
//ff:what pyFilenameFallback: attributes a Python function to its conventional test file via PyMatcher.Match when content-aware matching found nothing — the last-resort fallback the plan requires be preserved (Phase005b §2: the dynamic nature of Python, monkeypatch/fixtures, means static ast can miss calls, so test_<file>.py naming is a legitimate backup, not a hack). TestFuncs nil = run the whole file.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// pyFilenameFallback returns the conventional test_<file>.py (same dir or
// tests/) for fn via PyMatcher.Match, or found=false when none exists.
func pyFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	testRel, ok := (&PyMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
