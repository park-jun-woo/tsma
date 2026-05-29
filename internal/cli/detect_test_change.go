//ff:func feature=cli type=helper control=sequence
//ff:what Matches a function's tests and compares the combined mtime to detect changes
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// detectTestChange attributes the function's tests (content-aware for Go, file
// based otherwise) and compares the combined latest mtime of the matched test
// files against the stored fn.TestMtime. It returns (changed, match). When no
// test is attributed, it returns (false, empty match). Any one matched file
// being newer than the stored mtime makes changed=true.
func detectTestChange(root, lang string, fn *model.Function) (bool, match.TestMatch) {
	fm := match.NewFuncMatcher(lang)
	tm, found := fm.MatchFunc(root, fn)
	if !found || len(tm.Files) == 0 {
		return false, match.TestMatch{}
	}

	mtime := combinedTestMtime(root, tm.Files)
	if fn.TestMtime == mtime {
		return false, tm
	}
	return true, tm
}
