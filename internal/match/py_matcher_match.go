//ff:func feature=match type=implementation control=sequence
//ff:what Checks same dir and tests/ for test_ prefixed Python files matching the source name
package match

import (
	"os"
	"path/filepath"
)

func (m *PyMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	srcDir := filepath.Dir(sourceFile)
	srcBase := filepath.Base(sourceFile)
	expectedTest := "test_" + srcBase

	searchDirs := []string{srcDir, filepath.Join(srcDir, "tests")}

	for _, dir := range searchDirs {
		testRel := filepath.Join(dir, expectedTest)
		absPath := filepath.Join(projectRoot, testRel)
		if _, err := os.Stat(absPath); err == nil {
			return testRel, true
		}
	}

	return "", false
}
