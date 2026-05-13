//ff:func feature=match type=implementation control=sequence
//ff:what Checks same dir and __tests__/ for .test or .spec test files matching the source name
package match

import "path/filepath"

func (m *TSMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	srcDir := filepath.Dir(sourceFile)
	srcName := stripTSExtension(filepath.Base(sourceFile))

	searchDirs := []string{srcDir, filepath.Join(srcDir, "__tests__")}

	for _, dir := range searchDirs {
		if testRel, ok := findTSTestInDir(projectRoot, dir, srcName); ok {
			return testRel, true
		}
	}

	return "", false
}
