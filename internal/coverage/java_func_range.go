//ff:type feature=coverage type=model
//ff:what Helper struct for Java coverage analysis holding function location and name
package coverage

// javaFuncRange is a helper struct for Java coverage analysis.
type javaFuncRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
