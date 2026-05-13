//ff:func feature=match type=implementation control=iteration dimension=1
//ff:what Searches same dir and __tests__/ for .test/.spec files matching by filename or content
package match

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

func (m *TSMatcher) Match(projectRoot string, fn *model.Function) (string, bool) {
	srcDir := filepath.Join(projectRoot, filepath.Dir(fn.File))
	srcBase := filepath.Base(fn.File)
	srcName := stripTSExtension(srcBase)

	// Directories to search: same dir, then __tests__/ subdirectory.
	searchDirs := []string{srcDir, filepath.Join(srcDir, "__tests__")}

	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		if found, path := matchTSInDir(projectRoot, dir, entries, srcName, fn.Name); found {
			return path, true
		}
	}

	return "", false
}
