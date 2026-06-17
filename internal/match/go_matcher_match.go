//ff:func feature=match type=implementation control=sequence
//ff:what Checks if a corresponding _test.go file exists in the same directory
package match

import (
	"os"
	"path/filepath"
)

func (m *GoMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	// CanonicalTestPath is the single source of the <base>_test.go formula; ""
	// means the source is not a derivable Go file (non-.go), so no match.
	testRel := CanonicalTestPath("go", sourceFile)
	if testRel == "" {
		return "", false
	}
	absPath := filepath.Join(projectRoot, testRel)

	if _, err := os.Stat(absPath); err != nil {
		return "", false
	}

	return testRel, true
}
