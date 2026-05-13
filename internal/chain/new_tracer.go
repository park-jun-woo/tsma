//ff:func feature=chain type=factory control=selection
//ff:what Returns the appropriate Tracer implementation for a given language string
package chain

// NewTracer returns the appropriate tracer for the given language.
func NewTracer(lang string) Tracer {
	switch lang {
	case "go":
		return &GoTracer{}
	case "typescript":
		return &TSTracer{}
	case "python":
		return &PyTracer{}
	default:
		return &UnsupportedTracer{Lang: lang}
	}
}
