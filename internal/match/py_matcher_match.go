//ff:func feature=match type=implementation control=iteration dimension=1
//ff:what Searches same dir and tests/ for test_ prefix files matching by filename or content
package match

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

func (m *PyMatcher) Match(projectRoot string, fn *model.Function) (string, bool) {
	srcDir := filepath.Join(projectRoot, filepath.Dir(fn.File))
	srcBase := filepath.Base(fn.File)

	// Directories to search: same dir, then tests/ subdirectory.
	searchDirs := []string{srcDir, filepath.Join(srcDir, "tests")}

	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		if found, path := matchPyInDir(projectRoot, dir, entries, srcBase, fn.Name); found {
			return path, true
		}
	}

	return "", false
}
