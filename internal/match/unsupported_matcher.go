//ff:type feature=match type=implementation
//ff:what Returns no match for unsupported languages
package match

import "github.com/park-jun-woo/tsma/internal/model"

type unsupportedMatcher struct{}

func (u *unsupportedMatcher) Match(_ string, _ *model.Function) (string, bool) {
	return "", false
}
