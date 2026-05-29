//ff:type feature=coverage type=model
//ff:what Represents a flattened JaCoCo source file with its full relative path and per-line counters
package coverage

// jacocoFile is a source file flattened from the package/sourcefile nesting,
// carrying the full package-qualified relative path (e.g.
// "com/example/Calculator.java") alongside its per-line counters.
type jacocoFile struct {
	Path  string
	Lines []jacocoLine
}
