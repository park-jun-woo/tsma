//ff:func feature=graph type=helper control=sequence
//ff:what Returns true if the name starts with Test and has more characters
package graph

// isGoTestFunc returns true if the name starts with "Test" and has more characters.
func isGoTestFunc(name string) bool {
	if len(name) <= 4 {
		return false
	}
	return name[:4] == "Test"
}
