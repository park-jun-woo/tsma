//ff:type feature=chain type=model
//ff:what Stores metadata for a Python function definition including body lines
package chain

// pyFuncInfo stores metadata for a Python function definition.
type pyFuncInfo struct {
	name      string
	file      string // relative to project root
	startLine int
	endLine   int
	bodyLines []string // lines of the function body
}
