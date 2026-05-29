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
// matched function it returns its deduplicated TestMatch keyed by the function's
// index in fns; unmatched functions are absent from the result map. Only Go
// functions are content-aware here; callers handle other languages via the
// per-function fallback FuncMatcher.
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
		idx, err := BuildPkgTestIndex(projectRoot, dir)
		if err != nil {
			continue
		}
		for _, i := range idxs {
			refs, ok := MatchFuncByName(idx, &fns[i])
			if !ok {
				continue
			}
			tm, ok := refsToTestMatch(refs)
			if !ok {
				continue
			}
			out[i] = tm
		}
	}
	return out
}
