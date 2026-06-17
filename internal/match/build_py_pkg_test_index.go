//ff:func feature=match type=implementation control=iteration dimension=1 lang=python
//ff:what BuildPyPkgTestIndex: ast-parses every test_*.py / *_test.py in the package dir (and its tests/) — via ingestPyDir per directory — and builds the name→test-file index. The Python analogue of BuildTSPkgTestIndex. Returns nil when no Python interpreter is present or nothing was indexed, signaling MatchFunc to fall back to filename matching (PyMatcher).
package match

import "path/filepath"

// BuildPyPkgTestIndex builds the content-aware test index for pkgDir (relative
// to projectRoot), scanning pkgDir and pkgDir/tests. It returns nil when no
// Python interpreter is available or no test file yielded any reference.
func BuildPyPkgTestIndex(projectRoot, pkgDir string) *PyPkgTestIndex {
	python := resolvePython()
	if python == "" {
		return nil
	}
	idx := &PyPkgTestIndex{refs: make(map[string][]string)}
	for _, dir := range []string{pkgDir, filepath.Join(pkgDir, "tests")} {
		ingestPyDir(idx, python, projectRoot, dir)
	}
	if len(idx.refs) == 0 {
		return nil
	}
	return idx
}
