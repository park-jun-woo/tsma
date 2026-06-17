//ff:func feature=match type=helper control=sequence lang=rust
//ff:what rsFilenameFallback: attributes a Rust function to its conventional test file via RsMatcher.Match (the source file itself when it has an in-file #[cfg(test)] module, else tests/<name>.rs) when content-aware matching found nothing — the last-resort fallback the plan requires be preserved (not removed). TestFuncs nil = run the whole test binary.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// rsFilenameFallback returns the conventional test file for fn via
// RsMatcher.Match, or found=false when none exists.
func rsFilenameFallback(projectRoot string, fn *model.Function) (TestMatch, bool) {
	testRel, ok := (&RsMatcher{}).Match(projectRoot, fn.File)
	if !ok {
		return TestMatch{}, false
	}
	return TestMatch{Files: []string{testRel}, TestFuncs: nil}, true
}
