//ff:func feature=endpoint type=factory control=selection
//ff:what Returns the appropriate detector for the given language and framework
package endpoint

// NewDetector returns the appropriate detector for the given language/framework.
func NewDetector(lang, framework string) Detector {
	switch lang {
	case "go":
		switch framework {
		case "echo":
			return &GoEchoDetector{}
		case "chi":
			return &GoChiDetector{}
		default:
			return &GoGinDetector{}
		}
	case "typescript":
		switch framework {
		case "nextjs":
			return &TSNextjsDetector{}
		default:
			return &TSExpressDetector{}
		}
	case "python":
		switch framework {
		case "django":
			return &PyDjangoDetector{}
		default:
			return &PyFastapiDetector{}
		}
	default:
		return &UnsupportedDetector{Lang: lang}
	}
}
