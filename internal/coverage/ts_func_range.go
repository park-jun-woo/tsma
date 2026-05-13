//ff:type feature=coverage type=model
//ff:what Helper struct for TS/JS coverage analysis holding function location and name
package coverage

// tsFuncRange is a helper struct for coverage analysis of TypeScript functions.
type tsFuncRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
