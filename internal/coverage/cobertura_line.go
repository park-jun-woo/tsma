//ff:type feature=coverage type=model lang=csharp
//ff:what Represents a <line> element of a Cobertura XML export with hit count and branch condition coverage
package coverage

// coberturaLine is a <line> element. Number is the 1-based source line, Hits is
// the execution count, Branch reports whether the line is a branch point, and
// ConditionCoverage carries the "covered% (covered/total)" branch summary (e.g.
// "50% (1/2)") when Branch is "true".
type coberturaLine struct {
	Number            int    `xml:"number,attr"`
	Hits              int    `xml:"hits,attr"`
	Branch            string `xml:"branch,attr"`
	ConditionCoverage string `xml:"condition-coverage,attr"`
}
