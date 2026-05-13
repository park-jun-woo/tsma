//ff:func feature=graph type=factory control=selection
//ff:what Returns the appropriate Builder implementation for a given language string
package graph

// NewBuilder returns the appropriate builder for the given language.
func NewBuilder(lang string) Builder {
	switch lang {
	case "go":
		return &GoBuilder{}
	case "typescript":
		return &TSBuilder{}
	case "python":
		return &PyBuilder{}
	default:
		return &UnsupportedBuilder{Lang: lang}
	}
}
