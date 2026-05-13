//ff:func feature=runner type=factory control=selection
//ff:what Returns the appropriate Runner implementation for a given language string
package runner

// NewRunner returns the appropriate runner for the given language.
func NewRunner(lang string) Runner {
	switch lang {
	case "go":
		return &GoRunner{}
	case "typescript":
		return &TSRunner{}
	case "python":
		return &PyRunner{}
	default:
		return &UnsupportedRunner{Lang: lang}
	}
}
