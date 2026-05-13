//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Re-matches test files and groups functions by package directory
package cli

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// groupByPackage re-matches test files and groups functions by package directory.
func groupByPackage(root string, sess *model.Session, m match.Matcher) []pkgGroup {
	type key struct {
		pkgDir   string
		testFile string
	}
	idx := map[key]int{}
	var groups []pkgGroup

	for i := range sess.Functions {
		fn := &sess.Functions[i]
		testFile, found := m.Match(root, fn)
		fn.TestFile = testFile
		if !found {
			fn.Status = model.StatusTodo
			fn.FailOutput = ""
		}
		pkgDir := filepath.Dir(fn.File)
		k := key{pkgDir: pkgDir, testFile: testFile}
		pos, ok := idx[k]
		if !ok {
			idx[k] = len(groups)
			groups = append(groups, pkgGroup{pkgDir: pkgDir, testFile: testFile, funcs: []*model.Function{fn}})
		} else {
			groups[pos].funcs = append(groups[pos].funcs, fn)
		}
	}

	return groups
}
