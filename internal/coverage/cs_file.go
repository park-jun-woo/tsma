//ff:type feature=coverage type=model lang=csharp
//ff:what Represents a flattened Cobertura source file with its path and per-line counters
package coverage

// csFile is a source file flattened from the package/class nesting, carrying the
// class filename (e.g. "App/Calculator.cs") alongside its per-line counters.
type csFile struct {
	Path  string
	Lines []coberturaLine
}
