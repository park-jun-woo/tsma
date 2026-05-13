//ff:func feature=match type=implementation control=sequence
//ff:what Checks if a corresponding _test.go file exists in the same directory
package match

import (
	"os"
	"path/filepath"
	"strings"
)

func (m *GoMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	base := filepath.Base(sourceFile)
	if !strings.HasSuffix(base, ".go") {
		return "", false
	}

	testBase := strings.TrimSuffix(base, ".go") + "_test.go"
	testRel := filepath.Join(filepath.Dir(sourceFile), testBase)
	absPath := filepath.Join(projectRoot, testRel)

	if _, err := os.Stat(absPath); err != nil {
		return "", false
	}

	return testRel, true
}
