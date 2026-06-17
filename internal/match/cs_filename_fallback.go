//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what csFilenameFallback: attributes a C# function to its conventional FooTests.cs/FooTest.cs in the parallel *.Tests project via CsMatcher.Match when content-aware matching found nothing — the last-resort fallback the plan requires be preserved (not removed). TestFuncs nil = run the whole test class.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// csFilenameFallback returns the conventional FooTests.cs for fn via
// CsMatcher.Match, or found=false when none exists.
func csFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	testRel, ok := (&CsMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
