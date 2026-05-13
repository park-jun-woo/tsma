//ff:type feature=chain type=model
//ff:what Stores a discovered TypeScript/JavaScript function definition location
package chain

// tsFuncDef holds a discovered function definition in the project.
type tsFuncDef struct {
	name      string
	file      string // relative to project root
	startLine int
	endLine   int
}
