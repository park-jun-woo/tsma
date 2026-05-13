//ff:func feature=match type=factory control=selection
//ff:what Returns the appropriate Matcher implementation for a given language string
package match

// NewMatcher returns the appropriate Matcher for the given language.
func NewMatcher(lang string) Matcher {
	switch lang {
	case "go":
		return &GoMatcher{}
	case "typescript":
		return &TSMatcher{}
	case "python":
		return &PyMatcher{}
	default:
		return &unsupportedMatcher{}
	}
}
