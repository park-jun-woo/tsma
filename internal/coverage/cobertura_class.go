//ff:type feature=coverage type=model lang=csharp
//ff:what Represents a <class> element of a Cobertura XML export with its filename and per-line coverage
package coverage

// coberturaClass is a <class> element. Filename is the source path (relative to
// the report's <source> root, e.g. "App/Calculator.cs") and Lines holds the
// per-line hit and branch counters.
type coberturaClass struct {
	Name     string          `xml:"name,attr"`
	Filename string          `xml:"filename,attr"`
	Lines    []coberturaLine `xml:"lines>line"`
}
