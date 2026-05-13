//ff:type feature=coverage type=model
//ff:what Helper struct for Go coverage analysis holding function location and name
package coverage

// funcRange is a helper struct for coverage analysis.
type funcRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
