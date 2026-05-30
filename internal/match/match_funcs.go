//ff:func feature=match type=helper control=iteration dimension=2 lang=go
//ff:what Attributes tests to many Go functions, building each package index only once
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// MatchFuncs attributes tests to every function in fns using the content-aware
// Go index, building the index for each package directory exactly once and
// reusing it for all functions in that directory. This avoids the O(funcs ×
// files) re-parsing that calling MatchFunc per function would cause. For each
// function it first tries content-aware attribution (preserving precise 1:N
// multi-file matches); functions that content-aware does not match fall back to
// file-name matching via goFilenameFallback, which attributes the conventional
// <base>_test.go when it exists on disk. content-aware always wins: fallback is
// applied only to otherwise-unmatched functions and never overwrites a
// content-aware result. The result is keyed by the function's index in fns;
// functions with neither a content-aware match nor a conventional test file are
// absent from the map. This hybrid mirrors GoFuncMatcher.MatchFunc so analyze
// (batch) and detectTestChange (single) re-match identically. Only Go functions
// are content-aware here; callers handle other languages via the per-function
// fallback FuncMatcher.
func MatchFuncs(projectRoot string, fns []model.Function) map[int]TestMatch {
	out := make(map[int]TestMatch)

	// Group function indices by their package directory so each directory's
	// index is built once.
	byDir := make(map[string][]int)
	for i := range fns {
		dir := filepath.Dir(fns[i].File)
		byDir[dir] = append(byDir[dir], i)
	}

	for dir, idxs := range byDir {
		// idx and srcReceivers may be nil when the package cannot be parsed;
		// attributeFunc then relies solely on the file-name fallback (idx nil)
		// and treats every name as same-name-single (srcReceivers nil). Both
		// are built once per directory and reused for all funcs in it, mirroring
		// the single-func path's per-package build.
		idx, _ := BuildPkgTestIndex(projectRoot, dir)
		srcReceivers, _ := BuildPkgSourceReceivers(projectRoot, dir)
		for _, i := range idxs {
			if tm, ok := attributeFunc(projectRoot, idx, srcReceivers, &fns[i]); ok {
				out[i] = tm
			}
		}
	}
	return out
}
