//ff:func feature=runner type=util control=iteration dimension=1 lang=go
//ff:what Resolves the union of Go test function names for a match, extracting from files if absent
package runner

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// ResolveTestFuncs returns the set of Go test function names to run for a match.
// When the match already carries explicit TestFuncs (content-aware path), it
// returns them as-is. Otherwise (nil TestFuncs, e.g. a fallback or legacy
// single-file match) it parses every matched file and returns the deduplicated,
// order-preserving union of their TestXxx functions.
func ResolveTestFuncs(m match.TestMatch) ([]string, error) {
	if m.TestFuncs != nil {
		return m.TestFuncs, nil
	}
	seen := make(map[string]struct{})
	var funcs []string
	for _, f := range m.Files {
		abs, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}
		fns, err := ExtractTestFuncs(abs)
		if err != nil {
			return nil, err
		}
		funcs = appendUniqueFuncs(seen, funcs, fns)
	}
	return funcs, nil
}
