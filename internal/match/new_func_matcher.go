//ff:func feature=match type=factory control=selection
//ff:what Returns the appropriate FuncMatcher for a language (Go content-aware, others fallback)
package match

// NewFuncMatcher returns the appropriate FuncMatcher for the given language. Go
// and TypeScript get content-aware matchers; every other language (and the
// default) gets a fallback adapter that wraps the legacy file-name based
// Matcher, preserving existing behavior.
func NewFuncMatcher(lang string) FuncMatcher {
	switch lang {
	case "go":
		return &GoFuncMatcher{}
	case "typescript":
		return &TypeScriptFuncMatcher{}
	case "python":
		return &PythonFuncMatcher{}
	default:
		return &fallbackFuncMatcher{inner: NewMatcher(lang)}
	}
}
