//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what File-name fallback attribution for non-Go functions, preserving single-file behavior
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// attributeFallbackTests attributes tests to each non-Go function via the
// per-function file-name matcher, recording the single matched file onto its
// function. Behavior matches the legacy file-name matching path.
func attributeFallbackTests(projectRoot, lang string, functions []model.Function) {
	fm := match.NewFuncMatcher(lang)
	for i := range functions {
		tm, ok := fm.MatchFunc(projectRoot, &functions[i])
		if !ok {
			continue
		}
		setTestFiles(&functions[i], tm)
	}
}
