//ff:type feature=coverage type=model
//ff:what Holds the flattened set of JaCoCo source files parsed from a coverage report
package coverage

// jacocoCoverage is the parsed, flattened JaCoCo coverage: one entry per source
// file with its full relative path and per-line counters.
type jacocoCoverage struct {
	Files []jacocoFile
}
