//ff:func feature=match type=helper control=iteration dimension=1 lang=python
//ff:what ingestPyDir: scans one directory for test_*.py / *_test.py files and ingests each into the content index (via ingestPyTestFile). The Python analogue of ingestTSDir — extracted so the dir-loop and the entry-loop stay at separate depths. A missing directory is a no-op.
package match

import (
	"os"
	"path/filepath"
)

// ingestPyDir ingests every Python test file directly in dir (relative to
// projectRoot) into idx.
func ingestPyDir(idx *PyPkgTestIndex, python, projectRoot, dir string) {
	absDir := filepath.Join(projectRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !isPyTestFile(name) {
			continue
		}
		abs := filepath.Join(absDir, name)
		ingestPyTestFile(idx, python, abs, relTestPath(projectRoot, abs))
	}
}
