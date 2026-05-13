//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Iterates test suffixes to find a matching TS test file in a directory
package match

import (
	"os"
	"path/filepath"
)

// findTSTestInDir checks whether any test suffix variant exists in the given directory.
func findTSTestInDir(projectRoot, dir, srcName string) (string, bool) {
	for _, suffix := range tsTestSuffixes {
		testRel := filepath.Join(dir, srcName+suffix)
		absPath := filepath.Join(projectRoot, testRel)
		if _, err := os.Stat(absPath); err == nil {
			return testRel, true
		}
	}
	return "", false
}
