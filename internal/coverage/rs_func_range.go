//ff:type feature=coverage type=model
//ff:what Helper struct for Rust coverage analysis holding function location and name
package coverage

// rsFuncRange is a helper struct for Rust coverage analysis.
type rsFuncRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
