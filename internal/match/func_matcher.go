//ff:type feature=match type=factory
//ff:what Defines the FuncMatcher interface for content-aware function-to-test matching
package match

import "github.com/park-jun-woo/tsma/internal/model"

// FuncMatcher attributes tests to a source function. Unlike the file-name based
// Matcher, it works at function granularity and supports 1:N (one source
// function covered by several test functions/files).
type FuncMatcher interface {
	MatchFunc(projectRoot string, fn *model.Function) (TestMatch, bool)
}
