//ff:func feature=match type=helper control=sequence lang=java
//ff:what javaFilenameFallback: attributes a Java function to its conventional FooTest.java/FooTests.java in the src/test mirror via JavaMatcher.Match when content-aware matching found nothing — the last-resort fallback the plan requires be preserved (not removed). TestFuncs nil = run the whole test class.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// javaFilenameFallback returns the conventional FooTest.java for fn via
// JavaMatcher.Match, or found=false when none exists.
func javaFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	testRel, ok := (&JavaMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
