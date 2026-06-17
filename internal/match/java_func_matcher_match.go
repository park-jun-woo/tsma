//ff:func feature=match type=implementation control=sequence lang=java
//ff:what (JavaFuncMatcher).MatchFunc: builds the package's content-aware test index (scanning the JUnit src/test mirror dir) once and looks the function up by its bare name; on a content miss (no tree-sitter, no reference, or no test dir) it falls back to JavaMatcher filename matching. Mirrors TypeScriptFuncMatcher.MatchFunc so single-func re-matching stays content-precise without false TODOs.
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFunc attributes fn to the JUnit test files that call it (content-aware),
// or falls back to the conventional FooTest.java when nothing references it.
func (m *JavaFuncMatcher) MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool) {
	if fn == nil {
		return TestMatch{}, false
	}
	srcPkgDir := filepath.Dir(fn.File)
	idx := BuildJavaPkgTestIndex(projectRoot, srcPkgDir)
	if idx == nil {
		return javaFilenameFallback(projectRoot, fn)
	}
	if tm, ok := javaRefsToTestMatch(idx.refs[fn.Name]); ok {
		return tm, true
	}
	return javaFilenameFallback(projectRoot, fn)
}
