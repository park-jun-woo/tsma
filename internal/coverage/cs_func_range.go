//ff:type feature=coverage type=model lang=csharp
//ff:what Helper struct for C# coverage analysis holding function location and name
package coverage

// csFuncRange is a helper struct for C# coverage analysis.
type csFuncRange struct {
	file      string
	startLine int
	endLine   int
	funcName  string
}
