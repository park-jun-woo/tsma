//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Scans a directory for Python test files matching by filename then by content
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// matchPyInDir checks entries in a single directory for a Python test file match.
func matchPyInDir(projectRoot, dir string, entries []os.DirEntry, srcBase, funcName string) (bool, string) {
	expectedTest := "test_" + srcBase

	// Phase 1: filename-based match.
	for _, e := range entries {
		if e.IsDir() || e.Name() != expectedTest {
			continue
		}
		absPath := filepath.Join(dir, e.Name())
		rel, err := filepath.Rel(projectRoot, absPath)
		if err != nil {
			return true, absPath
		}
		return true, rel
	}

	// Phase 2: content-based match.
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasPrefix(e.Name(), "test_") || !strings.HasSuffix(e.Name(), ".py") {
			continue
		}
		absPath := filepath.Join(dir, e.Name())
		if !pyTestMentionsFunc(absPath, funcName) {
			continue
		}
		rel, err := filepath.Rel(projectRoot, absPath)
		if err != nil {
			return true, absPath
		}
		return true, rel
	}

	return false, ""
}
