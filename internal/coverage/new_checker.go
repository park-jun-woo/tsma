//ff:func feature=coverage type=factory control=selection
//ff:what Returns the appropriate Checker implementation for a given language string
package coverage

// NewChecker returns the appropriate checker for the given language.
func NewChecker(lang string) Checker {
	switch lang {
	case "go":
		return &GoChecker{}
	case "typescript":
		return &TSChecker{}
	case "python":
		return &PyChecker{}
	default:
		return &UnsupportedChecker{Lang: lang}
	}
}
