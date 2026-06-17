//ff:func feature=match type=implementation control=sequence lang=typescript
//ff:what (TypeScriptFuncMatcher).MatchFunc: builds the package's content-aware test index once and looks the function up by its bare name; on a content miss (no tree-sitter, no reference, or no files) it falls back to filename matching. Mirrors GoFuncMatcher.MatchFunc so single-func re-matching stays content-precise without false TODOs.
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc attributes fn to the test files that call it (content-aware), or
// falls back to the conventional sibling test file when nothing references it.
func (m *TypeScriptFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	pkgDir := filepath.Dir(fn.File)
	idx := BuildTSPkgTestIndex(projectRoot, pkgDir)
	if idx == nil {
		return tsFilenameFallback(projectRoot, fn)
	}
	if tm, ok := tsRefsToTestMatch(idx.refs[fn.Name]); ok {
		return tm, true
	}
	return tsFilenameFallback(projectRoot, fn)
}
