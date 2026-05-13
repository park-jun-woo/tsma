//ff:func feature=match type=implementation control=iteration dimension=1
//ff:what Scans *_test.go in the function's directory for Test* that contain the function name
package match

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

func (m *GoMatcher) Match(projectRoot string, fn *model.Function) (string, bool) {
	funcDir := filepath.Join(projectRoot, filepath.Dir(fn.File))

	entries, err := os.ReadDir(funcDir)
	if err != nil {
		return "", false
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		absPath := filepath.Join(funcDir, e.Name())
		if !containsTestFor(absPath, fn.Name) {
			continue
		}
		rel, err := filepath.Rel(projectRoot, absPath)
		if err != nil {
			return absPath, true
		}
		return rel, true
	}

	return "", false
}
