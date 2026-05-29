//ff:func feature=cli type=helper control=iteration dimension=1 lang=go
//ff:what Content-aware attribution for all Go functions, building each package index once
package cli

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// attributeGoTests attributes tests to every Go function using the content-aware
// batch matcher (one index build per package) and records each match onto its
// function.
func attributeGoTests(projectRoot string, functions []model.Function) {
	matches := match.MatchFuncs(projectRoot, functions)
	for i := range functions {
		tm, ok := matches[i]
		if !ok {
			continue
		}
		setTestFiles(&functions[i], tm)
	}
}
