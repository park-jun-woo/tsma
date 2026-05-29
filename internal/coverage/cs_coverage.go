//ff:type feature=coverage type=model lang=csharp
//ff:what Holds the flattened set of Cobertura source files parsed from a coverage report
package coverage

// csCoverage is the parsed, flattened Cobertura coverage: one entry per source
// file with its source-relative path and per-line counters.
type csCoverage struct {
	Files []csFile
}
