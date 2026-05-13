//ff:func feature=graph type=helper control=iteration dimension=1
//ff:what Updates brace depth by counting opening and closing braces in a line
package graph

// updateBraceDepth counts braces in a line and returns the updated depth.
func updateBraceDepth(line string, depth int) int {
	for _, ch := range line {
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
		}
	}
	return depth
}
