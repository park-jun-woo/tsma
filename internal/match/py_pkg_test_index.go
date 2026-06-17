//ff:type feature=match type=model lang=python
//ff:what PyPkgTestIndex: maps bare source identifier names to the Python test files that reference them — the content-aware index for one package directory. The Python analogue of TSPkgTestIndex: file granularity (pytest runs whole files), so TestFuncs is left nil downstream.
package match

// PyPkgTestIndex is a content-aware index for a single Python package directory.
// Keys are bare identifier names (function / method / class names) referenced by
// the directory's test files; each maps to the test files that reference it.
type PyPkgTestIndex struct {
	refs map[string][]string
}
