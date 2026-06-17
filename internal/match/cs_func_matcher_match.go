//ff:func feature=match type=implementation control=sequence lang=csharp
//ff:what (CsFuncMatcher).MatchFunc: builds the project's content-aware test index (scanning the parallel *.Tests project dirs) once and looks the function up by its bare name; on a content miss (no tree-sitter, no reference, or no test dir) it falls back to CsMatcher filename matching. Mirrors JavaFuncMatcher.MatchFunc so single-func re-matching stays content-precise without false TODOs.
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc attributes fn to the test files that call it (content-aware), or
// falls back to the conventional FooTests.cs when nothing references it.
func (m *CsFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	srcPkgDir := filepath.Dir(fn.File)
	idx := BuildCsPkgTestIndex(projectRoot, srcPkgDir)
	if idx == nil {
		return csFilenameFallback(projectRoot, fn)
	}
	if tm, ok := csRefsToTestMatch(idx.refs[fn.Name]); ok {
		return tm, true
	}
	return csFilenameFallback(projectRoot, fn)
}
