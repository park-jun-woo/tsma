//ff:type feature=coverage type=model
//ff:what Helper struct for Python coverage analysis holding function location and name
package coverage

// pyFuncRange is a helper struct for coverage analysis.
type pyFuncRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
