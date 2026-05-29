//ff:func feature=match type=implementation control=iteration dimension=1 lang=go
//ff:what Builds a content-aware test index for all _test.go files in a package dir
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// BuildPkgTestIndex parses every _test.go file in the package directory pkgDir
// (relative to projectRoot) and builds a content-aware index mapping bare
// source identifier names to the test functions that reference them. The index
// is built once per package and may be reused for every model.Function in that
// package. Unparseable test files are skipped. pkgDir "" or "." means the root.
func BuildPkgTestIndex(projectRoot, pkgDir string) (*PkgTestIndex, error) {
	idx := &PkgTestIndex{refs: make(map[string][]testRef)}
	absDir := filepath.Join(projectRoot, pkgDir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		abs := filepath.Join(absDir, name)
		ingestTestFile(idx, abs, relTestPath(projectRoot, abs))
	}
	return idx, nil
}
