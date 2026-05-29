//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Counts the net change in brace depth for a line, ignoring braces in strings, chars, and line comments
package index

// countBraces returns the net brace delta ('{' minus '}') for a single line,
// ignoring braces that appear inside string literals, char literals, or after
// a line comment. Block comments are not handled (best-effort).
func countBraces(line string) int {
	delta := 0
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == '/' && i+1 < len(line) && line[i+1] == '/':
			return delta // rest of line is a comment
		case c == '"' || c == '\'':
			i = skipQuoted(line, i)
		case c == '{':
			delta++
		case c == '}':
			delta--
		}
	}
	return delta
}
