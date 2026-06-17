//ff:func feature=match type=implementation control=sequence lang=rust
//ff:what (RsFuncMatcher).MatchFunc: builds the source file's content-aware test index (in-file #[cfg(test)] module + tests/*.rs) once and looks the function up by its bare name; on a content miss (no tree-sitter, no reference) it falls back to RsMatcher filename matching. Mirrors CsFuncMatcher.MatchFunc so single-func re-matching stays content-precise — a non-pub function the integration tests cannot see is still attributed via its in-file module.
package match

import "github.com/park-jun-woo/tsma/internal/model"

// MatchFunc attributes fn to the test files that call it (content-aware), or
// falls back to the conventional in-file / tests/<name>.rs when nothing
// references it.
func (m *RsFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	idx := BuildRsTestIndex(projectRoot, fn.File)
	if idx == nil {
		return rsFilenameFallback(projectRoot, fn)
	}
	if tm, ok := rsRefsToTestMatch(idx.refs[fn.Name]); ok {
		return tm, true
	}
	return rsFilenameFallback(projectRoot, fn)
}
