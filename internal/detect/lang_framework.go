//ff:type feature=detect type=model
//ff:what Holds language and framework detection results
package detect

// LangFramework holds detection results.
type LangFramework struct {
	Lang      string
	Framework string
}
