//ff:func feature=endpoint type=helper control=iteration dimension=1
//ff:what Extracts the leading whitespace from a line
package endpoint

// leadingWhitespace extracts the leading whitespace from a line.
func leadingWhitespace(line string) string {
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			return line[:i]
		}
	}
	return line
}
