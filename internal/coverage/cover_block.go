//ff:type feature=coverage type=model
//ff:what Represents a single coverage block from a Go coverprofile
package coverage

// coverBlock represents a single coverage block from a coverprofile.
type coverBlock struct {
	file      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
	stmts     int
	count     int
}
