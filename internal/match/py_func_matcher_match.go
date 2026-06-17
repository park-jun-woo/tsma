//ff:func feature=match type=implementation control=sequence lang=python
//ff:what (PythonFuncMatcher).MatchFunc: builds the package's content-aware test index once and looks the function up by its bare name; on a content miss (no interpreter, no reference, or no files) it falls back to filename matching (PyMatcher). Mirrors TypeScriptFuncMatcher.MatchFunc so single-func re-matching stays content-precise without false TODOs.
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc attributes fn to the test files that reference it (content-aware),
// or falls back to the conventional test_<file>.py when nothing references it.
func (m *PythonFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	pkgDir := filepath.Dir(fn.File)
	idx := BuildPyPkgTestIndex(projectRoot, pkgDir)
	if idx == nil {
		return pyFilenameFallback(projectRoot, fn)
	}
	if tm, ok := pyRefsToTestMatch(idx.refs[fn.Name]); ok {
		return tm, true
	}
	return pyFilenameFallback(projectRoot, fn)
}
