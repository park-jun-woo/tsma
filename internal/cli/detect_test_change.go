//ff:func feature=cli type=helper control=sequence
//ff:what Matches test file and compares mtime to detect changes
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// detectTestChange matches the test file and compares mtime.
// Returns (changed bool, testFile string). testFile is "" if no test file exists.
func detectTestChange(root, lang string, fn *model.Function) (bool, string) {
	m := match.NewMatcher(lang)
	testFile, found := m.Match(root, fn.File)
	if !found {
		return false, ""
	}

	mtime := getTestMtime(root, testFile)
	if fn.TestMtime == mtime {
		return false, testFile
	}

	return true, testFile
}
