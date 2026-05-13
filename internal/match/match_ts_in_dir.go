//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Scans a directory for TS test files matching by filename then by content
package match

import (
	"os"
	"path/filepath"
)

// matchTSInDir checks entries in a single directory for a TS test file match.
// It first tries filename-based matching, then content-based matching.
func matchTSInDir(projectRoot, dir string, entries []os.DirEntry, srcName, funcName string) (bool, string) {
	// Phase 1: filename-based match.
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if found, path := matchTSByFilename(projectRoot, dir, e.Name(), srcName); found {
			return true, path
		}
	}

	// Phase 2: content-based match.
	for _, e := range entries {
		if e.IsDir() || !isTSTestFile(e.Name()) {
			continue
		}
		absPath := filepath.Join(dir, e.Name())
		if !tsTestMentionsFunc(absPath, funcName) {
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
