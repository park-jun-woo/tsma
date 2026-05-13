//ff:type feature=match type=implementation
//ff:what Returns no match for unsupported languages
package match

type unsupportedMatcher struct{}

func (u *unsupportedMatcher) Match(_ string, _ string) (string, bool) {
	return "", false
}
